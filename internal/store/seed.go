package store

import (
	"github.com/gitvam/platform-go-challenge/internal/models"
)

func SeedDummyData(s *InMemoryStore) {
	// For user "johnsmith"

	// Chart
	s.AddFavorite("johnsmith", &models.Chart{
		ID:          "chart_engagement_2024",
		Title:       "Q1 2024 Social Media Engagement",
		XAxisTitle:  "Month",
		YAxisTitle:  "Engagement (k)",
		Data:        []int{85, 92, 110, 130},
		Description: "Tracks monthly engagement for all channels in Q1 2024.",
	})

	// Insight
	s.AddFavorite("johnsmith", &models.Insight{
		ID:          "insight_active_users",
		Text:        "78% of millennials engage with branded content daily.",
		Description: "Based on 2024 survey data across EMEA.",
	})

	// Audience
	s.AddFavorite("johnsmith", &models.Audience{
		ID:                 "aud_greece_men_24_35",
		Gender:             "male",
		BirthCountry:       "Greece",
		AgeGroups:          []string{"24-35"},
		HoursOnSocial:      4,
		PurchasesLastMonth: 3,
		Description:        "Digitally active Greek men aged 24-35 with high purchasing intent.",
	})

	// For user "maria.papadopoulou"

	// Chart
	s.AddFavorite("maria.papadopoulou", &models.Chart{
		ID:          "chart_ecom_conversion",
		Title:       "E-commerce Conversion Rates 2024",
		XAxisTitle:  "Week",
		YAxisTitle:  "Conversion Rate (%)",
		Data:        []int{2, 2, 3, 4, 3, 5, 4},
		Description: "Weekly conversion rate trend for Q2 2024.",
	})

	// Insight
	s.AddFavorite("maria.papadopoulou", &models.Insight{
		ID:          "insight_genz_tiktok",
		Text:        "Gen Z users are 3x more likely to purchase after seeing a TikTok ad.",
		Description: "Finding from global digital consumer study 2024.",
	})

	// Audience
	s.AddFavorite("maria.papadopoulou", &models.Audience{
		ID:                 "aud_uk_females_18_24",
		Gender:             "female",
		BirthCountry:       "UK",
		AgeGroups:          []string{"18-24"},
		HoursOnSocial:      6,
		PurchasesLastMonth: 5,
		Description:        "UK-based young women, highly active on Instagram and TikTok.",
	})
}
