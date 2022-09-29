package route

//import (
//	"io"
//	"mime/multipart"
//	"os"
//)
//
////定义 文件上传结构体，封装对文件的保存等操作
//
//type MultipartFile struct {
//	File map[string][]*multipart.FileHeader
//}
//
//// SaveUploadedFile 保存文件
//func (m *MultipartFile) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
//	src, err := file.Open()
//	if err != nil {
//		return err
//	}
//	defer src.Close()
//
//	out, err := os.Create(dst)
//	if err != nil {
//		return err
//	}
//	defer out.Close()
//
//	_, err = io.Copy(out, src)
//	return err
//}
