def main():
    TOKEN = '8178485363:AAGzYzstr-C6Gj9A8sR2MguA70f5wPFg6Q0'
    PORT = 8443  # Cổng nội bộ của bot, không cần phải thay đổi

    # Tạo đối tượng Updater
    updater = Updater(TOKEN, use_context=True)
    dp = updater.dispatcher

    # Tải proxy khi bot khởi động
    load_proxies()

    # Đăng ký các lệnh
    dp.add_handler(CommandHandler("on", start_bot))
    dp.add_handler(CommandHandler("off", stop_bot))
    dp.add_handler(CommandHandler("tang", tang))
    dp.add_handler(MessageHandler(Filters.document, update_proxy))

    # Khởi chạy webhook
    updater.start_webhook(
        listen="0.0.0.0",
        port=int(PORT),
        url_path=TOKEN,
        webhook_url=f"https://your-app-name.onrender.com/{TOKEN}"  # Đổi your-app-name thành tên ứng dụng Render của bạn
    )
    updater.idle()

if __name__ == '__main__':
    main()
