import http from "k6/http";
import { sleep } from "k6";

export let options = {
    scenarios: {
        start_flow: {           // часть юзеров просто запускают бота
            executor: "constant-vus",
            vus: 200,
            duration: "30s",
            exec: "startScenario",
        },
        booking_flow: {         // другие идут по сценарию записи
            executor: "constant-vus",
            vus: 750,
            duration: "30s",
            exec: "bookingScenario",
        },
        // myAppointments_flow: {  // кто-то смотрит свои записи
        //     executor: "constant-vus",
        //     vus: 30,
        //     duration: "30s",
        //     exec: "appointmentsScenario",
        // },
    },
};

// сценарий 1: /start
export function startScenario() {
    sendUpdate("/start");
    sleep(3); // пауза между действиями
}

// сценарий 2: запись на услугу
export function bookingScenario() {
    sendUpdate("/start");
    sleep(1);
    sendUpdate("Все услуги");
    sleep(1);
    sendUpdate("Парикмахер");
    sleep(1);
    sendUpdate("Дата: завтра");
    sleep(1);
    sendUpdate("10:00");
    sleep(1);
    sendUpdate("Иван");
    sleep(1);
    sendUpdate("1234567890");
    sleep(1);
    sendUpdate("Подтвердить");
    sleep(3);
}

// сценарий 3: просмотр записей
// export function appointmentsScenario() {
//     sendUpdate("/myappointments");
//     sleep(2);
//     sendUpdate("Имя: Иван");
//     sleep(2);
//     sendUpdate("Номер телефона: 1234567890");
//     sleep(3);
// }

// хелпер
function sendUpdate(text) {
    const payload = JSON.stringify({
        update_id: Math.floor(Math.random() * 1000000),
        message: {
            message_id: Math.floor(Math.random() * 10000),
            from: { id: Math.floor(Math.random() * 1000000), is_bot: false, first_name: "Test" },
            chat: { id: Math.floor(Math.random() * 1000000), type: "private" },
            date: Date.now(),
            text: text,
        },
    });

    http.post("http://go_backend:8080/fakeUpdate", payload, {
        headers: { "Content-Type": "application/json" },
    });
}
