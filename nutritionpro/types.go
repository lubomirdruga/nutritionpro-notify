package nutritionpro

import "time"

type MenuResponse struct {
	ID                 string    `json:"id"`
	Status             string    `json:"status"`
	Name               string    `json:"name"`
	User               User      `json:"user"`
	CurrentWeekDay     int       `json:"currentWeekDay"`
	Days               []Day     `json:"days"`
	StartDate          time.Time `json:"startDate"`
	EndDate            time.Time `json:"endDate"`
	Length             int       `json:"length"`
	Size               string    `json:"size"`
	MealsPerDay        int       `json:"mealsPerDay"`
	ClassMenu          string    `json:"classMenu"`
	Address            string    `json:"address"`
	DeliveryTime       []int     `json:"deliveryTime"`
	Energy             int       `json:"energy"`
	IgnoredMealTypes   []string  `json:"ignoredMealTypes"`
	AdditionalProducts []string  `json:"additionalProducts"`
	Polygon            Polygon   `json:"polygon"`
}

type User struct {
	ID                string      `json:"id"`
	InBodyID          int         `json:"inBodyId"`
	Name              string      `json:"name"`
	Phone             string      `json:"phone"`
	Address           string      `json:"address"`
	Email             string      `json:"email"`
	Language          string      `json:"language"`
	Country           string      `json:"country"`
	Comment           string      `json:"comment"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`
	SubjectID         string      `json:"subjectId"`
	BitrixID          int         `json:"bitrixId"`
	PaymentData       PaymentData `json:"paymentData"`
	LastDeal          LastDeal    `json:"lastDeal"`
	Birthday          time.Time   `json:"birthday"`
	ShowWizard        bool        `json:"showWizard"`
	Kitchen           string      `json:"kitchen"`
	Deleted           bool        `json:"deleted"`
	IsNotCzech        bool        `json:"isNotCzech"`
	IsIgnoreInRating  bool        `json:"isIgnoreInRating"`
	CreditBalance     int         `json:"creditBalance"`
	AffiliateID       *string     `json:"affiliateId"`
	IsReviewEvaluated bool        `json:"isReviewEvaluated"`
}

type PaymentData struct {
	IsCompany   bool   `json:"isCompany"`
	CompanyName string `json:"companyName"`
	Address     string `json:"address"`
	Zip         int    `json:"zip"`
	RegNumber   int    `json:"regNumber"`
	VatNumber   string `json:"vatNumber"`
}

type LastDeal struct {
	ID            int       `json:"id"`
	Date          time.Time `json:"date"`
	StageID       string    `json:"stageId"`
	StageName     string    `json:"stageName"`
	Reason        string    `json:"reason"`
	ReasonComment string    `json:"reasonComment"`
}

type Day struct {
	Timestamp int64     `json:"timestamp"`
	Dishes    []Dish    `json:"dishes"`
	Nutrients Nutrients `json:"nutrients"`
}

type Dish struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Meal        int       `json:"meal"`
	Weight      int       `json:"weight"`
	Size        float64   `json:"size"`
	IsHot       bool      `json:"isHot"`
	Nutrients   Nutrients `json:"nutrients"`
	Score       int       `json:"score"`
	Review      string    `json:"review"`
	DmuID       string    `json:"dmuId"`
	IsChoiced   bool      `json:"isChoiced"`
}

type Nutrients struct {
	Kcal float64 `json:"kcal"`
	Prot float64 `json:"prot"`
	Fat  float64 `json:"fat"`
	Carb float64 `json:"carb"`
}

type Polygon struct {
	Rating      int    `json:"rating"`
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	CountryCode string `json:"countryCode"`
}
