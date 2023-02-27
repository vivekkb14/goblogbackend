package httpserver_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vivekkb14/goblogbackend/common"
)

var _ = Describe("Dockerhub", func() {
	defer GinkgoRecover()

	Context("Tests for API operations of blogging server using Ginkgo", Ordered, func() {
		It("Test for displaying all the information", func() {
			fmt.Println("Hello world")
			var data common.DispalyAllAuthor
			url := "http://127.0.01:8080/articles"
			output, err := HttpReqest(http.MethodGet, url, nil)
			json.Unmarshal(output, &data)
			Expect(data.Status).To(Equal(200))
			Expect(err).ToNot(HaveOccurred())
		})

		It("Test case for insertion", func() {
			var data common.DispalyAllAuthor
			url := "http://127.0.01:8080/articles"
			payload := strings.NewReader(`{
				"title" : "Blogging Server",
				"content" : "This is server for blogging activities built using Golang",
				"author" : "VIVEK KB"
		}
		`)
			output, err := HttpReqest(http.MethodPost, url, payload)
			json.Unmarshal(output, &data)
			Expect(data.Status).To(Equal(200))
			Expect(err).ToNot(HaveOccurred())

		})

		It("Test case for insertion of invalid input", func() {
			var data common.DispalyAllAuthor
			url := "http://127.0.01:8080/articles"
			payload := strings.NewReader(`{
				"title" : "Blogging Server",
				"content" : "This is server for blogging activities built using Golang",
				"author" : 14
		}
		`)
			output, err := HttpReqest(http.MethodPost, url, payload)
			json.Unmarshal(output, &data)
			Expect(data.Status).To(Equal(500))
			Expect(err).ToNot(HaveOccurred())
		})

		It("Test for displaying present information in DB", func() {
			var data common.DispalyAllAuthor
			url := "http://127.0.01:8080/articles?id=1"
			output, err := HttpReqest(http.MethodGet, url, nil)
			json.Unmarshal(output, &data)
			Expect(data.Status).To(Equal(200))
			Expect(err).ToNot(HaveOccurred())
		})

		It("Test for wrong input to DB", func() {
			var data common.DispalyAllAuthor
			url := "http://127.0.01:8080/articles?id=100"
			output, err := HttpReqest(http.MethodGet, url, nil)
			Expect(err).ToNot(HaveOccurred())
			json.Unmarshal(output, &data)
			Expect(data.Status).To(Equal(http.StatusInternalServerError))
		})

	})
})

func HttpReqest(method, url string, payload io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(body))
	return body, nil
}
