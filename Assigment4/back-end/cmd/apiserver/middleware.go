package main

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"taskmanager/models"
	"time"
)

// responseWriter is a custom ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set security headers
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		//logger.WithField("path", r.URL.Path).Info("Applied security headers")

		next.ServeHTTP(w, r)
	})
}

func TeamRoleMiddleware(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value("user_id").(uint)
			if !ok {
				logger.Warn("user_id missing in context")
				http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
				return
			}

			teamID, err := strconv.Atoi(mux.Vars(r)["team_id"])
			if err != nil {
				logger.WithError(err).Warn("Invalid team_id in request")
				http.Error(w, `{"error": "Invalid team_id"}`, http.StatusBadRequest)
				return
			}

			// Check user's role in the team
			hasRole, err := CheckTeamRole(uint(teamID), userID, requiredRoles...)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"user_id": userID,
					"team_id": teamID,
				}).WithError(err).Warn("User role check failed")
				http.Error(w, `{"error": "Unauthorized or not part of the team"}`, http.StatusUnauthorized)
				return
			}

			if !hasRole {
				logger.WithFields(logrus.Fields{
					"user_id":        userID,
					"team_id":        teamID,
					"required_roles": requiredRoles,
				}).Warn("Insufficient permissions")
				http.Error(w, `{"error": "Insufficient permissions"}`, http.StatusForbidden)
				return
			}

			logger.WithFields(logrus.Fields{
				"user_id":       userID,
				"team_id":       teamID,
				"roles_checked": requiredRoles,
			}).Info("User authorized for team action")

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func RequestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		logger.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"status":   ww.statusCode,
			"duration": time.Since(start).Milliseconds(),
		}).Info("Request completed")

	})
}

func ModeMiddleware(nextIndividual, nextTeam http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			logger.Warn("user_id missing in context")
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Fetch user from the database
		var user models.User
		if err := DB.First(&user, userID).Error; err != nil {
			logger.WithField("user_id", userID).WithError(err).Warn("User not found")
			http.Error(w, `{"error": "User not found"}`, http.StatusUnauthorized)
			return
		}

		// Direct request based on mode
		if user.Mode == "individual" {
			logger.WithField("user_id", userID).Info("User operating in individual mode")
			nextIndividual(w, r)
		} else if user.Mode == "team" {
			logger.WithField("user_id", userID).Info("User operating in team mode")
			nextTeam(w, r)
		} else {
			logger.WithField("user_id", userID).Warn("Invalid user mode")
			http.Error(w, `{"error": "Invalid user mode"}`, http.StatusBadRequest)
		}
	}
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()

		// Record metrics
		requestCount.WithLabelValues(r.Method, r.URL.Path, http.StatusText(ww.statusCode)).Inc()
		responseDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}

func ExpvarMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		totalRequests.Add(1)

		// Call the next handler
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)

		// Count errors
		if ww.statusCode >= 400 {
			errorRequests.Add(1)
		}
	})
}
