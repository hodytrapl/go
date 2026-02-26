package model

import(
	"time"
)

// result - дата хранения данных о ссылке
type Result struct {
	URL 		string 			`json:"url"` 						//ссылка
	Status 		int 			`json:"Status"` 					//статус - 404,200
	OK 			bool 			`json:"ok"` 						//ответ запроса - ок:200, неок:404
	TTFB 		time.Duration 	`json:"-"` 							//таймаут - время на проверку
	TTFBms 		int64 			`json:"ttfb_ms"` 					//таймаут в милисекундах
	SizeBytes 	int64 			`json:"size_bytes"` 				//размер тела
	Match 		bool 			`json:"match"` 						//
	Location 	string 			`json:"Location, omitempty"` 		//
	Error 		string 			`json:"error, omitempty"` 			//ошибка
}

func(r *Result) Finalize(){
	r.TTFBms=r.TTFB.Milliseconds()
}