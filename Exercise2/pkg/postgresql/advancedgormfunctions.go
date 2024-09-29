package postgresql

import (
	"exercise2/models"
	"fmt"
	"gorm.io/gorm"
)

func InsertUserWithProfile(db *gorm.DB) error {
	// Transaction block
	return db.Transaction(func(tx *gorm.DB) error {
		// Create the User first
		user := models.AdvancedGormUser{
			Name: "Alice",
			Age:  30,
		}

		// Insert the User record
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("error inserting user: %v", err)
		}

		fmt.Printf("Inserted user with ID: %d\n", user.ID)

		// Explicitly set the UserID in the Profile and insert the Profile
		profile := models.AdvancedGormProfile{
			AdvancedGormUserID: user.ID, // Assign the created user's ID
			Bio:                "Software Developer",
			ProfilePictureURL:  "http://example.com/alice.jpg",
		}

		// Insert the Profile record
		if err := tx.Create(&profile).Error; err != nil {
			return fmt.Errorf("error inserting profile: %v", err)
		}

		fmt.Println("User and Profile inserted successfully!")
		return nil
	})
}

func GetUsersWithProfiles(db *gorm.DB) ([]models.AdvancedGormUser, error) {
	var users []models.AdvancedGormUser

	// Eager load Profile with Preload
	if err := db.Preload("Profile").Find(&users).Error; err != nil {
		return nil, err
	}

	// Display the results
	for _, user := range users {
		fmt.Printf("User: %s, Age: %d, Bio: %s, Profile Picture: %s\n",
			user.Name, user.Age, user.Profile.Bio, user.Profile.ProfilePictureURL)
	}

	return users, nil
}

func UpdateUserProfile(db *gorm.DB, userID uint, newBio string, newProfilePicURL string) error {
	// Update the profile of the user with the given ID
	err := db.Model(&models.AdvancedGormProfile{}).Where("advanced_gorm_user_id = ?", userID).
		Updates(models.AdvancedGormProfile{Bio: newBio, ProfilePictureURL: newProfilePicURL}).Error
	if err != nil {
		return err
	}
	fmt.Println("User profile updated successfully!")
	return nil
}

func DeleteUserWithProfile(db *gorm.DB, userID uint) error {
	// Delete the user, associated Profile will be deleted due to foreign key constraint
	err := db.Delete(&models.AdvancedGormUser{}, userID).Error
	if err != nil {
		return err
	}
	fmt.Println("User and associated Profile deleted successfully!")
	return nil
}
