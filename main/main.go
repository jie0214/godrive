package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// serviceAccount
// 可使用環境變數GOOGLE_APPLICATION_CREDENTIALS指定檔案路徑
const serviceAccountFilePath string = "key/poc-center-googledrive.json"

func main() {

	ctx := context.Background()

	srv, err := drive.NewService(ctx, option.WithServiceAccountFile(serviceAccountFilePath))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	//如果文件量過多可以調整PageSize
	r, err := srv.Files.List().PageSize(10).
		//nextPageToken 透過Toke換頁
		Fields("nextPageToken, files(id, name, mimeType, parents)").
		Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	for _, f := range r.Files {
		//跳過資料夾
		if f.MimeType == "application/vnd.google-apps.folder" {
			continue
		}
		if err := Download(ctx, srv, f.Name, f.Id); err != nil {
			log.Fatalf("Unable to download: %v", err)
		}
	}
}

func Download(ctx context.Context, srv *drive.Service, name, id string) error {
	// create, err := os.Create(name)
	// if err != nil {
	// 	return fmt.Errorf("create file: %w", err)
	// }
	// defer create.Close()

	// resp, err := srv.Files.Get(id).Context(ctx).Download()
	// if err != nil {
	// 	return fmt.Errorf("get drive file: %w", err)
	// }
	// defer resp.Body.Close()

	// if _, err := io.Copy(create, resp.Body); err != nil {
	// 	return fmt.Errorf("write file: %w", err)
	// }

	// MIME sheet to csv
	mimeType := "text/csv"
	//創建文件
	create, err := os.Create(name + ".csv")
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer create.Close()
	//繪出內容
	resp, err := srv.Files.Export(id, mimeType).Download()
	if err != nil {
		return fmt.Errorf("get drive file: %w", err)
	}
	if _, err := io.Copy(create, resp.Body); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func GrafanaSecheduleUpdate() {

}
