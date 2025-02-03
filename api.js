const axios = require("axios");

// Danh sách API (tự động loại bỏ "/api" nếu có)
const API_URL = [
    "https://dos-ime7.onrender.com/api",
    "https://dos-tls.onrender.com/api",
    "https://dos-ura.onrender.com/api",
    "https://dosbuoi.onrender.com/api",
    "https://dosflood.onrender.com/api",
    "https://quangdevclo.onrender.com/api"
];

// Hàm gọi API tự động mỗi 40 giây
const autoCallAPI = async () => {
    for (const api of API_URL) {
        const cleanApi = api.replace(/\/api$/, ""); // Loại bỏ /api nếu có
        try {
            await axios.get(cleanApi);
            console.log(`[✔] Gọi API thành công: ${cleanApi}`);
        } catch (error) {
            console.log(`[❌] Gọi API thất bại: ${cleanApi} - Lỗi: ${error.message}`);
        }
    }
};

// Thiết lập gọi API mỗi 40 giây
setInterval(autoCallAPI, 40 * 1000);

console.log("🟢 Chương trình tự động gọi API đang chạy...");
