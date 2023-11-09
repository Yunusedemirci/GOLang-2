package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/common-nighthawk/go-figure"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Pink   = "\033[35m"
)

func main() {

	myFigure := figure.NewFigure("YAVUZLAR", "", true)
	myFigure.Print()
	myFigure2 := figure.NewFigure("WEB SCAPPER TOOL", "small", true)
	myFigure2.Print()

	websiteFlag1 := flag.Bool("1", false, "TheHackerNews Haber Sitesinden Veri Çek")
	//flag.Bool fonksiyonu, terminalden girilen parametrelerin boolean (true/false) olacağını belirtir.
	websiteFlag2 := flag.Bool("2", false, "SonDakika Haber Sitesinden Veri Çek")
	dateFlag := flag.Bool("date", false, "Tarih bilgilerini görüntülemek istemiyorsanız kullanın")
	descFlag := flag.Bool("description", false, "Açıklama bilgilerini görüntülemek istemiyorsanız kullanın")

	flag.Parse()
	//flag.Parse() fonksiyonu ile terminalden girilen parametreleri alıyoruz.

	var websites []string
	//websites adında bir slice oluşturuyoruz.
	//slice, dizilerden farklı olarak boyutu dinamik olarak değişebilen veri tipleridir.
	//slice şu işe yarar mesela, 10 elemanlı bir dizi oluşturduk ve 10 elemanı da doldurduk.
	//ama 11. elemanı eklemek istiyoruz, bunu dizilerde yapamayız ama slice ile yapabiliriz.
	//slice'ın boyutunu arttırabiliriz.

	if *websiteFlag1 {
		websites = append(websites, "http://www.thehackernews.com/")
		//append fonksiyonu ile slice'ın sonuna eleman ekleyebiliriz.
		//append fonksiyonu, slice'ın boyutunu arttırır.
		//append fonksiyonu, slice'ın sonuna eklenen elemanı döndürür.

	}

	if *websiteFlag2 {
		websites = append(websites, "https://www.sondakika.com/teknoloji/")
	}

	for _, url := range websites {
		// _, url ifadesindeki _ işareti, Go dilinde kullanılmayan değişkenleri temsil eder
		// yani _ işareti, url değişkeninin değerini almıyor.
		// range websites ifadesi, websites slice'ının elemanlarını tek tek döndürür.
		// for döngüsü ile slice'ın elemanlarını tek tek döndürüyoruz.
		res, err := http.Get(url)
		// res, err := http.Get(url) ifadesi ile url değişkenindeki adrese GET isteği gönderiyoruz.
		// http.Get fonksiyonu, http paketinin bir fonksiyonudur.

		if err != nil {
			// eğer bir hata varsa, err değişkeni nil olmayacaktır.
			// nil, Go dilinde null değerine karşılık gelir.
			// err değişkeni nil değilse, hata var demektir.
			log.Fatal("Hata:", err)
			return
		}
		defer res.Body.Close()
		// defer ifadesi, fonksiyonun sonunda çalıştırılacak kod bloğunu belirtir.
		// res.Body.Close() ifadesi ile http.Get fonksiyonu ile açılan bağlantıyı kapatıyoruz.

		if res.StatusCode != 200 {
			log.Fatal("Hata: Sayfa yüklenemedi, HTTP kodu:", res.StatusCode)
			return
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		// goquery.NewDocumentFromReader(res.Body) ifadesi ile http.Get fonksiyonu ile açılan bağlantıdan gelen veriyi okuyoruz.
		// goquery paketinin NewDocumentFromReader fonksiyonu, verilen okuyucudan (reader) veri okur.
		// res.Body değişkeni, http.Get fonksiyonu ile açılan bağlantıdan gelen veriyi tutar.
		// res.Body değişkeni, io.ReadCloser türündedir.
		// io.ReadCloser türü, io.Reader ve io.Closer türlerini içerir.
		// io.Reader türü, veri okuyucuları için bir arayüz (interface) türüdür.

		if err != nil {
			log.Fatal("Hata:", err)
			return
		}

		if strings.Contains(url, "thehackernews") {
			// strings.Contains fonksiyonu, bir string içerisinde başka bir string arar.
			// eğer aranan string bulunursa, true değerini döndürür.
			// eğer aranan string bulunmazsa, false değerini döndürür.
			doc.Find(".clear.home-right").Each(func(i int, s *goquery.Selection) {
				// doc.Find fonksiyonu, belirtilen sorgu ifadesine göre HTML elementlerini seçer.
				// doc.Find fonksiyonu, goquery.Document türünden bir değer döndürür.
				// goquery.Document türü, HTML elementlerini seçmek için kullanılır.
				// goquery.Document türü, goquery.Selection türünden bir değer döndürür.
				// goquery.Selection türü, seçilen HTML elementlerini tutar.
				// goquery.Selection türü, HTML elementlerini seçmek için kullanılır.

				title := s.Find(".home-title").Text()

				datetime := s.Find(".h-datetime").Text()
				desc := s.Find(".home-desc").Text()

				printData(title, datetime, desc, *dateFlag, *descFlag)
			})
		} else if strings.Contains(url, "sondakika") {
			doc.Find(".nws").Each(func(i int, s *goquery.Selection) {
				title := s.Find(".title").Text()
				datetime := s.Find(".hour.data_calc").Text()
				desc := s.Find(".news-detail.news-column").Text()

				printData(title, datetime, desc, *dateFlag, *descFlag)
			})
		}
	}
}

func printData(title, datetime, desc string, hideDate, hideDesc bool) {
	fmt.Println(Pink + "Başlık:\n" + Reset + title)

	if !hideDate {

		fmt.Println(Yellow + "Tarih:\n" + Reset + datetime)
	}

	if !hideDesc {
		fmt.Println(Green + "Açıklama:\n" + Reset + desc)
	}

	fmt.Println("")
	fmt.Println("--------------------------------------------------")
}
