import http from "k6/http"
import { sleep } from "k6"

// const TEST_DURATION = "10m"
// const PERIOD = "5m"
const MIN_RPS = 0;
const MAX_RPS = 100;

function stages() {
    const duration = 10 * 60 // min * sec = duration sec
    const period = 2 * 60 // min * sec = period sec
    let stages = []
    const steps = 60;

    for (let i = 0; i <= steps; i++) {
        const time = i * duration / steps
        const rps = MIN_RPS + (MAX_RPS - MIN_RPS) * (1 + Math.sin(2 * Math.PI * time / period)) / 2

        stages.push({
            duration: Math.round(duration / steps)+"s",
            target: Math.round(rps)
        })
    }
    // console.log(stages);

    return stages
}

export const options = {
    // stages: stages(),
    discardResponseBodies: true,
    scenarios: {
        sine_wave: {
            executor: 'ramping-arrival-rate',
            startRate: MIN_RPS,
            timeUnit: '1s',
            preAllocatedVUs: 20,
            maxVUs: 200,
            stages: stages().map(stage => ({
                target: stage.target,
                duration: stage.duration
            }))
        }
    }
}

export default function() {
    const payload = JSON.stringify({
        product_id: "15",
        count: "5",
        username: "jo.peach",
    })
    const headers = { 'Content-Type': 'application/json' };

    http.post("http://localhost:8083/create", payload, {headers})

    sleep(1)
}