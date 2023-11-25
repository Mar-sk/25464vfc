package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {
	coins := []string{
		"PAZUSDT", "FAVRUSDT", "MYRIAUSDT", "LAIUSDT", "AZYUSDT", "MMITUSDT", "WATTUSDT", "RBCUSDT", "XRPCUSDT", "LBRUSDT", "DOGSROCKUSDT", "DPXUSDT", "LIEFLATUSDT", "LIFEUSDT", "LSDUSDT", "PHBUSDT", "ACSUSDT", "ISKUSDT", "DLANCEUSDT", "EMTUSDT", "METISUSDT", "YOOSHIUSDT", "TENETUSDT", "HIUSDT", "CAPOUSDT", "OMMIUSDT", "AYAUSDT", "GUILDUSDT", "FIDAUSDT", "JPEGUSDT", "MXENUSDT", "CALCIUMUSDT", "OPOSUSDT", "LITHUSDT", "OPNUSDT", "BEAMUSDT", "SWASHUSDT", "MDIUSDT", "SQUID2USDT", "ZOONUSDT", "USDJUSDT", "ITHEUMUSDT", "IDYPUSDT", "THNUSDT", "VELODROMEUSDT", "RLYUSDT", "VISIONUSDT", "DGCUSDT", "PHAUSDT", "DODOUSDT", "PRQUSDT", "ZBCUSDT", "PEOPLEUSDT", "DRMUSDT", "RIWAUSDT", "JSTUSDT", "ROACHCOINUSDT", "AGIIUSDT", "LOOMUSDT", "AIUSDT", "LBCUSDT", "BEAIUSDT", "MXTUSDT", "FYNUSDT", "CHEUSDT", "AMEUSDT", "UNLEASHUSDT", "KASTAUSDT", "ROSXUSDT", "AICODEUSDT", "DAIUSDT", "HATIUSDT", "FREEDOMUSDT", "FINUUSDT", "ICTUSDT", "WILDUSDT", "AGROUSDT", "IONXUSDT", "GNUSDT", "EVMOSUSDT", "CGVUSDT", "POOHUSDT", "TRAVAUSDT", "BTCBAMUSDT", "BAXUSDT", "MNTUSDT", "XELSUSDT", "KZENUSDT", "PORTUMAUSDT", "FAKTUSDT", "AINUSDT", "CREUSDT", "JINKOUSDT", "KISHUUSDT", "FCT2USDT", "KAFUSDT", "MBLUSDT", "EDUUSDT", "MEFAUSDT", "ZCXUSDT"
	}
	// Бесконечный цикл для проверки каждую новую минуту
	for {
		found := false // Флаг для отслеживания наличия монет, удовлетворяющих условиям

		// Перебираем все монеты в списке
		for i, symbol := range coins {
			// URL для получения исторических свечей
			apiURL := fmt.Sprintf("https://api.mexc.com/api/v3/klines?symbol=%s&interval=1m&limit=2", symbol)

			// Выполняем GET-запрос к бирже MEXC
			response, err := http.Get(apiURL)
			if err != nil {
				fmt.Println("Ошибка при выполнении GET-запроса:", err)
				continue
			}
			defer response.Body.Close()

			// Читаем ответ API
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Ошибка при чтении ответа:", err)
				continue
			}
			fmt.Println(i)

			// Разбираем JSON-данные в массив массивов
			var klines [][]interface{}
			err = json.Unmarshal(responseBody, &klines)
			if err != nil {
				fmt.Println("Ошибка при разборе JSON:", err)
				continue
			}
			// Выводим данные свечей
			fmt.Printf("\rДанные свечей для %s:", symbol)
			for _, candle := range klines {
				fmt.Println("\r", candle)
			}

			if len(klines) < 2 {
				fmt.Println("Недостаточно данных для вычислений")
				continue
			}

			// Получаем последние две свечи
			x := klines[0]
			y := klines[1]

			// Получаем объем актива для двух свечей
			volumeX, err := strconv.ParseFloat(x[5].(string), 64)
			if err != nil {
				fmt.Println("Ошибка при разборе объема предпоследней свечи:", err)
				continue
			}

			volumeY, err := strconv.ParseFloat(y[5].(string), 64)
			if err != nil {
				fmt.Println("Ошибка при разборе объема последней свечи:", err)
				continue
			}

			// Расчет изменения объема и процента изменения
			volumeChange := volumeY - volumeX
			percentageChange := (volumeChange / volumeX) * 100.0

			// Получаем цены открытия и закрытия для последней свечи
			openPriceY, err := strconv.ParseFloat(y[1].(string), 64)
			if err != nil {
				fmt.Println("Ошибка при разборе цены открытия последней свечи:", err)
				continue
			}

			closePriceY, err := strconv.ParseFloat(y[4].(string), 64)
			if err != nil {
				fmt.Println("Ошибка при разборе цены закрытия последней свечи:", err)
				continue
			}

			// Расчет изменения цены открытия и закрытия
			priceChange := closePriceY - openPriceY
			// Расчет изменения цены в процентах
			priceChangePercentage := (priceChange / openPriceY) * 100.0

			// Проверяем условия
			if (priceChangePercentage >= 8.0) && ((percentageChange >= 300.0) || (volumeX == 0.0 && volumeY > 0.0)) {
				// Если оба условия выполняются, выводим сообщение "ПАМП"
				fmt.Printf("ПАМП %s Процент изменений цены: %.15f%%\n", symbol, priceChangePercentage)
				bot, _ := tgbotapi.NewBotAPI("6855593649:AAEijSEkGmvbTU8mjOIquxFco809O_ApW64")

				chatID := 2107725915
				message := fmt.Sprintf("ПАМП %s Процент изменений цены: %.15f%%\n", symbol, priceChangePercentage)
				msg := tgbotapi.NewMessage(int64(chatID), message)
				bot.Send(msg)

				found = true // Устанавливаем флаг, что монета удовлетворяет условиям
			}
		}

		// Если не найдено монет, удовлетворяющих условиям, выводим сообщение
		if !found {
			fmt.Println("Нет монет удовлетворяющих условиям")
		}

		// Засыпаем до наступления следующей минуты
		now := time.Now()
		nextMinute := now.Truncate(time.Minute).Add(time.Minute)
		time.Sleep(time.Until(nextMinute))
	}
}
