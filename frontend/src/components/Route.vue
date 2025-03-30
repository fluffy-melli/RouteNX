<script setup>
import { ref, onMounted } from 'vue'
import { Chart, LinearScale, CategoryScale, LineElement, PointElement, Title, Tooltip, Legend, LineController } from 'chart.js'

Chart.register(
  LinearScale,
  CategoryScale,
  LineElement,
  PointElement,
  Title,
  Tooltip,
  Legend,
  LineController
)

const data = ref([])

const debug = false
const api = debug ? 'http://localhost:3000' : window.location.origin

onMounted(async () => {
  const response = await fetch(`${api}/route`)
  if (response.ok) {
    data.value = await response.json()
  } else {
    console.error('Failed to fetch data')
  }

  let ctx = document.getElementById('traffic').getContext('2d')
  let chart = new Chart(ctx, {
    type: 'line',
    data: {
      labels: [],
      datasets: [
        {
          label: 'TX',
          borderColor: 'rgb(75, 192, 192)',
          data: [],
          fill: false,
          tension: 0.3
        },
        {
          label: 'RX',
          borderColor: 'rgb(255, 99, 132)',
          data: [],
          fill: false,
          tension: 0.3
        },
      ],
    },
    options: {
      responsive: true,
      aspectRatio: 1.5,
      maintainAspectRatio: false,
      animation: {
        duration: 0,
      },
      scales: {
        x: {
          type: 'linear',
          position: 'bottom',
          ticks: {
            callback: function(value) {
              const time = new Date(value)
              const min = time.getMinutes()
              const scd = time.getSeconds()
              return min.toString() +":"+ scd.toString()
            }
          }
        },
        y: {
          beginAtZero: true,
          ticks: {
            callback: function(value) {
              if (value >= 1073741824) { // 1024 * 1024 * 1024
                return (value / 1073741824).toFixed(2) + ' GBps'
              } else if (value >= 1048576) { // 1024 * 1024
                return (value / 1048576).toFixed(2) + ' MBps'
              } else if (value >= 1024) {
                return (value / 1024).toFixed(2) + ' KBps'
              } else {
                return value + ' Bps'
              }
            }
          }
        }
      },
    },
  })
  setInterval(() => {
    getTraffic(chart)
  }, 1500)
})

async function getTraffic(chart) {
  const response = await fetch(`${api}/traffic`)
  if (response.ok) {
    const trf = await response.json()
    console.log(trf)
    chart.data.labels.push(...trf["Label"])
    chart.data.datasets[0].data.push(...trf["TX"])
    chart.data.datasets[1].data.push(...trf["RX"])
    if (chart.data.labels.length > 20) {
      chart.data.labels = chart.data.labels.slice(-20);
      chart.data.datasets[0].data = chart.data.datasets[0].data.slice(-20);
      chart.data.datasets[1].data = chart.data.datasets[1].data.slice(-20);
    }
    chart.update()
  } else {
    console.error('Failed to fetch data')
  }
}

function getSettingClass(name, list, index) {
  const isFirst = index === 0
  const isLast = index === list.length - 1

  if (isFirst && list.length === 1) {
    return `${name}-all`
  } else if (isFirst && !isLast) {
    return `${name}-start`
  } else if (isLast && !isFirst) {
    return `${name}-end`
  } else {
    return `${name}-none`
  }
}

function formatEndpoint(endpoint) {
  if (endpoint.startsWith('http://') || endpoint.startsWith('https://')) {
    return endpoint.replace(/^https?:\/\//, '')
  }
  return endpoint
}
</script>

<template>
  <div class="app">
    <div>
      <div>
        <h2 class="title">RouteNX</h2>
        <div class="line"></div>
      </div>
    </div>
    <div class="container">
      <h2 class="title">Reverse Proxy</h2>
      <p class="description">List of available routes and proxy settings</p>
      <table class="styled-table">
        <thead>
          <tr>
            <th class="long">Hostname</th>
            <th class="short">Firewall</th>
            <th class="long">Endpoint</th>
            <th class="actions">Actions</th>
            <th></th>
          </tr>
        </thead>
        <tbody v-if="data.routes">
          <tr v-for="(item, index) in data.routes" :key="index">
            <td>
              <select>
                <option v-for="(item, index) in item.host" :key="index" :value="item" :disabled="index !== 0">
                  {{ item }}
                </option>
              </select>
            </td>
            <td>
              <select v-if="item.firewall.length > 0">
                <option v-for="(item, index) in item.firewall" :key="index" :value="item" :disabled="index !== 0">
                  {{ item }}
                </option>
              </select>
              <select v-else>
                <option>
                  x
                </option>
              </select>
            </td>
            <td>
              <select>
                <option>
                  {{ formatEndpoint(item.endpoint) }}
                </option>
              </select>
            </td>
            <td>
              <div class="setting">
                <img class="remove" src="../assets/svg/remove.svg" alt="remove"/>
              </div>
            </td>
            <td></td>
          </tr>
          <tr class="none">
            <td>
              <p class="lefts">+ Add more</p>
            </td>
            <td></td>
            <td></td>
            <td>
              <p class="actions">{{ data.routes.length }}</p>
            </td>
            <td></td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="line"></div>
    <div class="container">
      <h2 class="title">Firewall</h2>
      <p class="description">List of firewall rules</p>
      <table class="styled-table">
        <thead>
          <tr>
            <th class="long">Rulename</th>
            <th class="short">To</th>
            <th class="long">CIDR</th>
            <th class="actions">Actions</th>
            <th></th>
          </tr>
        </thead>
        <tbody v-if="data.firewall">
          <tr v-for="(item, index) in data.firewall" :key="index">
            <td>
              <select>
                <option>
                  {{ item.name }}
                </option>
              </select>
            </td>
            <td>
              <select>
                <option v-if="!item.block">
                  Allow
                </option>
                <option v-else>
                  Block
                </option>
              </select>
            </td>
            <td>
              <select>
                <option v-for="(item, index) in item.cidr" :key="index" :value="item" :disabled="index !== 0">
                  {{ item }}
                </option>
              </select>
            </td>
            <td>
              <div class="setting">
                <img class="remove" src="../assets/svg/remove.svg" alt="remove"/>
              </div>
            </td>
            <td></td>
          </tr>
          <tr class="none">
            <td>
              <p class="lefts">+ Add more</p>
            </td>
            <td></td>
            <td></td>
            <td>
              <p class="actions">{{ data.firewall.length }}</p>
            </td>
            <td></td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="line"></div>
    <div class="container">
      <canvas id="traffic"></canvas>
    </div>
  </div>
</template>