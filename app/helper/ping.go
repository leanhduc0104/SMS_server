package helper

import (
	"context"
	"os/exec"
	"time"
)

// PingServer is a function to ping a server
func PingServer(address string) error {
	// Thiết lập thời gian chờ là 3 giây
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	// Tạo lệnh ping
	cmd := exec.CommandContext(ctx, "ping", address)

	// Thực thi lệnh và lấy kết quả
	_, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return ctx.Err()
	}

	if err != nil {
		return err
	}

	return nil
}
