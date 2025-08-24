import http from "k6/http";
import { check } from "k6";

export let options = {
    vus: 100,         //  виртуальных пользователей
    duration: "10s", // тест n секунд
};

export default function () {
    const payload = JSON.stringify({
        update_id: Math.floor(Math.random() * 1000000),
        message: {
            message_id: Math.floor(Math.random() * 10000),
            from: { id: Math.floor(Math.random() * 1000000), is_bot: false, first_name: "Test" },
            chat: { id: Math.floor(Math.random() * 1000000), type: "private" },
            date: Date.now(),
            text: "Все услуги"
        }
    });

    const res = http.post("http://go_backend:8081/fakeUpdate", payload, {
        headers: { "Content-Type": "application/json" },
    });

    check(res, { "status is 200": (r) => r.status === 200 });
}
