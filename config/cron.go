package config

import (
	"fmt"
	"log"
	"time"

	"github.com/rickyromansyah2045/halocat-backend-go/content"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"github.com/rickyromansyah2045/halocat-backend-go/logs"
	"github.com/rickyromansyah2045/halocat-backend-go/user"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func InitScheduler(db *gorm.DB) {
	jakartaTime, err := time.LoadLocation("Asia/Jakarta")

	if err != nil {
		log.Fatal("error while load time location, err: ", err.Error())
	}

	scheduler := cron.New(cron.WithLocation(jakartaTime))

	defer scheduler.Stop()

	// for testing: */1 * * * *
	// for prod: 10 0 * * *
	scheduler.AddFunc("10 0 * * *", func() {
		affected := 0
		rows, err := db.Raw(helper.ConvertToInLineQuery(content.QueryGetAll+"AND status = 'active' AND finished_at <= ?"), fmt.Sprintf("%v", time.Now().Format("2006-01-02 15:04:05"))).Rows()
		activityLog := logs.ActivityLog{}
		activityLog.IpAddress = "-"
		activityLog.UserAgent = "-"

		if err != nil {
			activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 1)] %v", err.Error())

			log.Println(activityLog.Content)

			if err := db.Create(&activityLog).Error; err != nil {
				log.Fatal(err.Error())
			}
		}

		defer rows.Close()

		for rows.Next() {
			tmp := content.Content{}
			err := rows.Scan(
				&tmp.ID,
				&tmp.UserID,
				&tmp.CategoryID,
				&tmp.Title,
				&tmp.Slug,
				&tmp.ShortDescription,
				&tmp.Description,
				&tmp.Status,
				&tmp.FinishedAt,
				&tmp.CreatedAt,
				&tmp.CreatedBy,
				&tmp.UpdatedAt,
				&tmp.UpdatedBy,
				&tmp.DeletedAt,
				&tmp.DeletedBy,
			)

			if err != nil {
				activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 2)] %v", err.Error())

				log.Println(activityLog.Content)

				if err := db.Create(&activityLog).Error; err != nil {
					log.Fatal(err.Error())
				}
			}

			result := db.Model(&content.Content{}).Where("id = ?", tmp.ID).Update("status", "finished")

			if result.Error != nil {
				activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 4-1)] %v", err.Error())

				log.Println(activityLog.Content)

				if err := db.Create(&activityLog).Error; err != nil {
					log.Fatal(err.Error())
				}
			}

			var userData user.User

			if err := db.Where("id = ?", tmp.UserID).Find(&userData).Error; err != nil {
				activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 4-2)] %v", err.Error())

				log.Println(activityLog.Content)

				if err := db.Create(&activityLog).Error; err != nil {
					log.Fatal(err.Error())
				}
			}

			if userData.ID == 0 {
				activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 4-3)] %v", "sql: no rows in result set")

				log.Println(activityLog.Content)

				if err := db.Create(&activityLog).Error; err != nil {
					log.Fatal(err.Error())
				}
			}

			affected = affected + int(result.RowsAffected)
		}

		activityLog.Content = fmt.Sprintf("System running CRON for check and update finished campaign. (affected: %v)", affected)

		log.Println(activityLog.Content)

		if err := db.Create(&activityLog).Error; err != nil {
			log.Fatal(err.Error())
		}
	})

	go scheduler.Start()
}
