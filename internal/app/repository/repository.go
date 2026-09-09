package repository

import "fmt"

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusDeleted   = "deleted"
)

type DepositMonth struct {
	ID               int
	Name             string
	ShortDescription string
	Description      string
	MonthNumber      int
	DaysCount        int
	Status           string
	ImageURL         string
	VideoURL         string
	Likes []int
}

type Repository struct {
	months []DepositMonth
}

func NewRepository() *Repository {
	return &Repository{
		months: []DepositMonth{
			{
				ID:               1,
				Name:             "Январь",
				ShortDescription: "Начало годового цикла начислений.",
				Description:      "В январе теплее, когда деньги не просто лежат, а понемногу растут на вкладе.",
				MonthNumber:      1,
				DaysCount:        31,
				Status:           StatusPublished,
				ImageURL:         "http://localhost:9000/deposit-month-assets/january.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/january.mp4",
				Likes:            []int{1, 3, 5, 8},
			},
			{
				ID:               2,
				Name:             "Февраль",
				ShortDescription: "Завершение зимнего периода начислений.",
				Description:      "Февраль короткий, но даже за 28 дней вклад успевает принести прибыль.",
				MonthNumber:      2,
				DaysCount:        28,
				Status:           StatusPublished,
				ImageURL:         "http://localhost:9000/deposit-month-assets/february.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/february.mp4",
				Likes:            []int{2, 4, 6},
			},
			{
				ID:               3,
				Name:             "Март",
				ShortDescription: "Начало весеннего периода начислений.",
				Description:      "Весной распускаются первые листочки, а твои накопления продолжают расти",
				MonthNumber:      3,
				DaysCount:        31,
				Status:           StatusPublished,
				ImageURL:         "http://localhost:9000/deposit-month-assets/march.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/march.mp4",
				Likes:            []int{1, 2, 4, 7, 9},
			},
			{
				ID:               4,
				Name:             "Апрель",
				ShortDescription: "Весенний расчётный период.",
				Description:      "Пока за окном идут апрельские дожди, и проценты по вкладу продолжают капать.",
				MonthNumber:      4,
				DaysCount:        30,
				Status:           StatusPublished,
				ImageURL:         "http://localhost:9000/deposit-month-assets/april.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/april.mp4",
				Likes:            []int{2, 8},
			},
			{
				ID:               5,
				Name:             "Май",
				ShortDescription: "Завершение весеннего периода начислений.",
				Description:      "Май - это время цветения, накопления тоже потихоньку расцветают.",
				MonthNumber:      5,
				DaysCount:        31,
				Status:           StatusPublished,
				ImageURL:         "http://localhost:9000/deposit-month-assets/may.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/may.mp4",
				Likes:            []int{1, 3, 6, 7},
			},
			{
				ID:               6,
				Name:             "Июнь",
				ShortDescription: "Начало летнего периода начислений.",
				Description:      "Лето только начинается, а деньги уже могут работать, пока ты отдыхаешь.",
				MonthNumber:      6,
				DaysCount:        30,
				Status:           StatusPublished,
				ImageURL:         "http://localhost:9000/deposit-month-assets/june.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/june.mp4",
				Likes:            []int{5, 8, 10},
			},
			{
				ID:               7,
				Name:             "Июль",
				ShortDescription: "Летний расчётный период.",
				Description:      "Пока июль радует солнцем и лимонадом, вклад радует начисленными процентами.",
				MonthNumber:      7,
				DaysCount:        31,
				Status:           StatusPublished,
				ImageURL:         "http://localhost:9000/deposit-month-assets/july.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/july.mp4",
				Likes:            []int{1, 4, 6},
			},
			{
				ID:               8,
				Name:             "Август",
				ShortDescription: "Завершение летнего периода начислений.",
				Description:      "Август - это отличная возможность проверить, насколько успели подрасти накопления за лето.",
				MonthNumber:      8,
				DaysCount:        31,
				Status:           StatusPublished,
				ImageURL:         "http://localhost:9000/deposit-month-assets/august.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/august.mp4",
				Likes:            []int{2, 3},
			},

			{
				ID:               9,
				Name:             "Сентябрь",
				ShortDescription: "Начало осеннего периода начислений.",
				Description:      "Осень начинается спокойнее, если знаешь, что работаешь не только ты, но и деньги",
				MonthNumber:      9,
				DaysCount:        30,
				Status:           StatusDraft,
				ImageURL:         "http://localhost:9000/deposit-month-assets/september.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/september.mp4",
				Likes:            []int{},
			},

			{
				ID:               10,
				Name:             "Октябрь",
				ShortDescription: "Удалённая карточка.",
				Description:      "Карточка октября удалена.",
				MonthNumber:      10,
				DaysCount:        31,
				Status:           StatusDeleted,
				ImageURL:         "http://localhost:9000/deposit-month-assets/october.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/october.mp4",
				Likes:            []int{},
			},
			{
				ID:               11,
				Name:             "Ноябрь",
				ShortDescription: "Удалённая карточка.",
				Description:      "Карточка ноября удалена.",
				MonthNumber:      11,
				DaysCount:        30,
				Status:           StatusDeleted,
				ImageURL:         "http://localhost:9000/deposit-month-assets/november.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/november.mp4",
				Likes:            []int{},
			},
			{
				ID:               12,
				Name:             "Декабрь",
				ShortDescription: "Удалённая карточка.",
				Description:      "Карточка декабря удалена.",
				MonthNumber:      12,
				DaysCount:        31,
				Status:           StatusDeleted,
				ImageURL:         "http://localhost:9000/deposit-month-assets/december.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/december.mp4",
				Likes:            []int{},
			},
		},
	}
}

func (r *Repository) GetPublishedDepositMonths(days *int) []DepositMonth {
	var result []DepositMonth

	for _, month := range r.months {
		if month.Status != StatusPublished {
			continue
		}

		if days != nil && month.DaysCount != *days {
			continue
		}

		result = append(result, month)
	}

	return result
}

func (r *Repository) GetDepositMonth(id int) (DepositMonth, error) {
	for _, month := range r.months {
		if month.ID == id && month.Status == StatusPublished {
			return month, nil
		}
	}

	return DepositMonth{}, fmt.Errorf(
		"опубликованный месяц с id %d не найден",
		id,
	)
}

func (r *Repository) GetFirstPublishedDepositMonth() (DepositMonth, error) {
	for _, month := range r.months {
		if month.Status == StatusPublished {
			return month, nil
		}
	}

	return DepositMonth{}, fmt.Errorf("опубликованные месяцы не найдены")
}

func (r *Repository) GetDraftDepositMonth() (DepositMonth, error) {
	for _, month := range r.months {
		if month.Status == StatusDraft {
			return month, nil
		}
	}

	return DepositMonth{}, fmt.Errorf("черновик не найден")
}

func (r *Repository) GetNextPublishedDepositMonth(id int) (DepositMonth, error) {
	currentIndex := -1

	for i, month := range r.months {
		if month.ID == id && month.Status == StatusPublished {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return DepositMonth{}, fmt.Errorf(
			"опубликованный месяц с id %d не найден",
			id,
		)
	}

	for step := 1; step < len(r.months); step++ {
		index := (currentIndex + step) % len(r.months)

		if r.months[index].Status == StatusPublished {
			return r.months[index], nil
		}
	}

	return DepositMonth{}, fmt.Errorf("следующий месяц не найден")
}
