package repository

import "fmt"

const (
	Deposit_month_status_draft     = "draft"
	Deposit_month_status_published = "published"
	Deposit_month_status_deleted   = "deleted"
)

type Deposit_month struct {
	ID               int
	Name             string
	ShortDescription string
	Description      string

	MonthNumber int
	DaysCount   int

	Status string

	ImageURL string
	VideoURL string

	Likes []int
}

type Deposit_month_repository struct {
	deposit_months []Deposit_month
}

func New_deposit_month_repository() *Deposit_month_repository {
	return &Deposit_month_repository{
		deposit_months: []Deposit_month{
			{
				ID:               1,
				Name:             "Январь",
				ShortDescription: "Начало годового цикла начислений.",
				Description:      "В январе теплее, когда деньги не просто лежат, а понемногу растут на вкладе.",
				MonthNumber:      1,
				DaysCount:        31,
				Status:           Deposit_month_status_published,
				ImageURL:         "http://localhost:9000/deposit-month-assets/january.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/january.mp4",
				Likes:            []int{1, 3, 5, 8},
			},
			{
				ID:               2,
				Name:             "Февраль",
				ShortDescription: "Второй расчётный период года.",
				Description:      "Февраль короткий, но даже за 28 дней вклад успевает принести начисленные проценты.",
				MonthNumber:      2,
				DaysCount:        28,
				Status:           Deposit_month_status_published,
				ImageURL:         "http://localhost:9000/deposit-month-assets/february.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/february.mp4",
				Likes:            []int{2, 4, 6},
			},
			{
				ID:               3,
				Name:             "Март",
				ShortDescription: "Начало весеннего периода начислений.",
				Description:      "Весной растёт всё — пусть вместе с первыми листьями растут и накопления.",
				MonthNumber:      3,
				DaysCount:        31,
				Status:           Deposit_month_status_published,
				ImageURL:         "http://localhost:9000/deposit-month-assets/march.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/march.mp4",
				Likes:            []int{1, 2, 4, 7, 9},
			},
			{
				ID:               4,
				Name:             "Апрель",
				ShortDescription: "Весенний расчётный период.",
				Description:      "Пока за окном идут апрельские дожди, проценты по вкладу продолжают начисляться.",
				MonthNumber:      4,
				DaysCount:        30,
				Status:           Deposit_month_status_published,
				ImageURL:         "http://localhost:9000/deposit-month-assets/april.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/april.mp4",
				Likes:            []int{2, 8},
			},
			{
				ID:               5,
				Name:             "Май",
				ShortDescription: "Завершение весеннего периода начислений.",
				Description:      "Май — время цветения, а накопления тоже могут постепенно расти.",
				MonthNumber:      5,
				DaysCount:        31,
				Status:           Deposit_month_status_published,
				ImageURL:         "http://localhost:9000/deposit-month-assets/may.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/may.mp4",
				Likes:            []int{1, 3, 6, 7},
			},
			{
				ID:               6,
				Name:             "Июнь",
				ShortDescription: "Начало летнего периода начислений.",
				Description:      "Лето только начинается, а деньги уже могут работать, пока вы отдыхаете.",
				MonthNumber:      6,
				DaysCount:        30,
				Status:           Deposit_month_status_published,
				ImageURL:         "http://localhost:9000/deposit-month-assets/june.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/june.mp4",
				Likes:            []int{5, 8, 10},
			},
			{
				ID:               7,
				Name:             "Июль",
				ShortDescription: "Летний расчётный период.",
				Description:      "Пока июль радует солнцем, вклад продолжает приносить начисленные проценты.",
				MonthNumber:      7,
				DaysCount:        31,
				Status:           Deposit_month_status_published,
				ImageURL:         "http://localhost:9000/deposit-month-assets/july.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/july.mp4",
				Likes:            []int{1, 4, 6},
			},
			{
				ID:               8,
				Name:             "Август",
				ShortDescription: "Завершение летнего периода начислений.",
				Description:      "Август — хороший момент посмотреть, сколько успели подрасти накопления за лето.",
				MonthNumber:      8,
				DaysCount:        31,
				Status:           Deposit_month_status_published,
				ImageURL:         "http://localhost:9000/deposit-month-assets/august.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/august.mp4",
				Likes:            []int{2, 3},
			},
			{
				ID:               9,
				Name:             "Сентябрь",
				ShortDescription: "Начало осеннего периода начислений.",
				Description:      "Осень начинается спокойнее, когда накопления продолжают работать сами.",
				MonthNumber:      9,
				DaysCount:        30,
				Status:           Deposit_month_status_draft,
				ImageURL:         "http://localhost:9000/deposit-month-assets/september.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/september.mp4",
				Likes:            []int{},
			},
			{
				ID:               10,
				Name:             "Октябрь",
				ShortDescription: "Осенний расчётный период.",
				Description:      "Октябрьский период используется для расчёта процентов за 31 календарный день.",
				MonthNumber:      10,
				DaysCount:        31,
				Status:           Deposit_month_status_deleted,
				ImageURL:         "http://localhost:9000/deposit-month-assets/october.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/october.mp4",
				Likes:            []int{},
			},
			{
				ID:               11,
				Name:             "Ноябрь",
				ShortDescription: "Поздний осенний расчётный период.",
				Description:      "Ноябрьский период используется для расчёта процентов за 30 календарных дней.",
				MonthNumber:      11,
				DaysCount:        30,
				Status:           Deposit_month_status_deleted,
				ImageURL:         "http://localhost:9000/deposit-month-assets/november.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/november.mp4",
				Likes:            []int{},
			},
			{
				ID:               12,
				Name:             "Декабрь",
				ShortDescription: "Завершение годового цикла начислений.",
				Description:      "Декабрь завершает годовой цикл расчёта ежемесячных процентных выплат.",
				MonthNumber:      12,
				DaysCount:        31,
				Status:           Deposit_month_status_deleted,
				ImageURL:         "http://localhost:9000/deposit-month-assets/december.jpg",
				VideoURL:         "http://localhost:9000/deposit-month-assets/december.mp4",
				Likes:            []int{},
			},
		},
	}
}

func (r *Deposit_month_repository) Get_published_deposit_months(
	min_days_count int,
	max_days_count int,
) []Deposit_month {

	result := make([]Deposit_month, 0)

	for _, deposit_month := range r.deposit_months {

		if deposit_month.Status != Deposit_month_status_published {
			continue
		}

		if deposit_month.DaysCount < min_days_count {
			continue
		}

		if deposit_month.DaysCount > max_days_count {
			continue
		}

		result = append(
			result,
			deposit_month,
		)
	}

	return result
}

func (r *Deposit_month_repository) Get_deposit_month(
	deposit_month_id int,
) (Deposit_month, error) {

	for _, deposit_month := range r.deposit_months {

		if deposit_month.ID == deposit_month_id &&
			deposit_month.Status == Deposit_month_status_published {

			return deposit_month, nil
		}
	}

	return Deposit_month{},
		fmt.Errorf(
			"deposit_month с id %d не найден",
			deposit_month_id,
		)
}

func (r *Deposit_month_repository) Get_first_published_deposit_month() (
	Deposit_month,
	error,
) {

	for _, deposit_month := range r.deposit_months {

		if deposit_month.Status ==
			Deposit_month_status_published {

			return deposit_month, nil
		}
	}

	return Deposit_month{},
		fmt.Errorf(
			"опубликованные deposit_month не найдены",
		)
}

func (r *Deposit_month_repository) Get_draft_deposit_month() (
	Deposit_month,
	error,
) {

	for _, deposit_month := range r.deposit_months {

		if deposit_month.Status ==
			Deposit_month_status_draft {

			return deposit_month, nil
		}
	}

	return Deposit_month{},
		fmt.Errorf(
			"draft deposit_month не найден",
		)
}

func (r *Deposit_month_repository) Get_next_published_deposit_month(
	deposit_month_id int,
) (Deposit_month, error) {

	current_index := -1

	for index, deposit_month := range r.deposit_months {

		if deposit_month.ID == deposit_month_id &&
			deposit_month.Status == Deposit_month_status_published {

			current_index = index
			break
		}
	}

	if current_index == -1 {
		return Deposit_month{},
			fmt.Errorf(
				"deposit_month с id %d не найден",
				deposit_month_id,
			)
	}

	for step := 1; step < len(r.deposit_months); step++ {

		index :=
			(current_index + step) %
				len(r.deposit_months)

		if r.deposit_months[index].Status ==
			Deposit_month_status_published {

			return r.deposit_months[index], nil
		}
	}

	return Deposit_month{},
		fmt.Errorf(
			"следующий deposit_month не найден",
		)
}
