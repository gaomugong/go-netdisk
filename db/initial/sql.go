package initial

import (
	"go-netdisk/db"
	"go-netdisk/db/models"
	"go-netdisk/settings"
	"go-netdisk/utils"
	"log"
)

func InitData() {
	if settings.ENV.NeedMigrate {
		_ = db.DB.AutoMigrate(&models.Project{}, &models.User{}, &models.Permission{}, &models.Matter{}, models.Preference{})
	}

	log.Printf("Create superuser: %s", settings.ENV.SuperUser)
	if _, err := models.GetOrCreateUser(settings.ENV.SuperUser, true); err != nil {
		panic(err)
	}

	perm := &models.Permission{}
	db.DB.Where(models.Permission{UserName: settings.ENV.SuperUser}).Attrs(models.Permission{
		Role: models.ADMINISTRATOR,
	}).FirstOrCreate(&perm)
	log.Printf("GetOrCreate permission: %s\n", utils.PrettyJson(perm))

	prefer := &models.Preference{}
	db.DB.Where(models.Preference{Name: "netdisk"}).Attrs(models.Preference{
		AllowRegister: true,
	}).FirstOrCreate(&prefer)
	log.Printf("GetOrCreate preference: %s\n", utils.PrettyJson(prefer))

	if db.DB.First(&models.Project{}).RowsAffected == 0 {
		log.Printf("Create default project")
		project := models.Project{
			Name:        "DEMO",
			Description: "DEMO",
		}
		if err := db.DB.Create(&project).Error; err != nil {
			panic(err)
		}

		guest, _ := models.GetOrCreateUser("user0", false)
		db.DB.Create(&models.Permission{
			UserName:    guest.Username,
			ProjectUUID: project.UUID,
			Role:        "USER",
		})
	}

}
