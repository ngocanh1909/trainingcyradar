## trainingcyradar ##
**task1**
  1. Request đến https://malshare.com/daily
  2. Đọc Response (Response sẽ chính là text HTML - bản chất là string)
  3. Duyệt string để lấy các thông tin cần thiết (đường link đến các ngày) --> Có thể sử dụng regex để việc lấy các thông tin cần thiết đơn giản.
  4. Sau khi lấy được link đến từng ngày thì lại Request đến link đó (có dạng https://malshare.com/daily/yyyy-MM-dd/malshare_fileList.yyyy-MM-dd.all.txt)
  Nhận Response là các md5, sha1, sha256, ....
  5. Lưu các thông tin cần có lại thành 1 map: ngày --> md5, ngày --> sha1, ....
  6. Duyệt map vừa lưu, tạo các thư mục tương ứng theo ngày, tháng, năm và viết file.
