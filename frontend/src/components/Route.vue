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
const selected = ref({})

const debug = true
const api = debug ? 'http://localhost:3000' : window.location.origin

const baseProxyRule = {
  "host": [
    ""
  ],
  "firewall": [
    ""
  ],
  "endpoint": ""
}

const baseFireWallRule = {
  "name": "",
  "To": false,
  "cidr": [
    ""
  ]
}

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
      chart.data.labels = chart.data.labels.slice(-20)
      chart.data.datasets[0].data = chart.data.datasets[0].data.slice(-20)
      chart.data.datasets[1].data = chart.data.datasets[1].data.slice(-20)
    }
    chart.update()
  } else {
    console.error('Failed to fetch data')
  }
}

function saveConfig() {
  fetch(`${api}/config`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data.value),
  })
  .then(response => {
      if (!response.ok) {
        throw new Error('Failed to save config')
      }
      return response.json()
  })
  .then(async (response) => {
    data.value = response
    alert('saved')
  })
  .catch(error => {
    console.error('Error:', error)
    alert('Failed')
  })
}


function formatEndpoint(endpoint) {
  if (endpoint.startsWith('http://') || endpoint.startsWith('https://')) {
    return endpoint.replace(/^https?:\/\//, '')
  }
  return endpoint
}

function selectRule(name, index) {
  if (selected.value[name] != null && selected.value[name] === index) {
    selected.value = { ...selected.value, [name]: null }
  } else {
    selected.value = { ...selected.value, [name]: index }
  }
}

function removeRule(name, index) {
  if (this.selected[name] == index) {
    selected.value = { ...selected.value, [name]: null }
  }
  this.data[name].splice(index, 1)
}

function addRule(name) {
  if (name == 'routes') {
    this.data[name].push(baseProxyRule)
  } else {
    this.data[name].push(baseFireWallRule)
  }
  this.selected[name] = this.data[name].length - 1
}

function removeItem(name, type, index) {
  this.data[name][this.selected[name]][type].splice(index, 1)
}

function addItem(name, type) {
  this.data[name][this.selected[name]][type].push("")
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
            <th class="short">Actions</th>
          </tr>
        </thead>
        <tbody v-if="data.routes">
          <tr class="rule" v-for="(item, index) in data.routes" :key="index" @click="selectRule('routes', index)">
            <td>
              <select class="select" @click.stop>
                <option v-for="(item, index) in item.host" :key="index" :value="item" :disabled="index !== 0">
                  {{ item }}
                </option>
              </select>
            </td>
            <td>
              <select class="select" v-if="item.firewall.length > 0" @click.stop>
                <option v-for="(item, index) in item.firewall" :key="index" :value="item" :disabled="index !== 0">
                  {{ item }}
                </option>
              </select>
              <select class="select" v-else @click.stop>
                <option>
                  x
                </option>
              </select>
            </td>
            <td>
              <select class="select" @click.stop>
                <option>
                  {{ formatEndpoint(item.endpoint) }}
                </option>
              </select>
            </td>
            <td>
              <div class="setting">
                <img class="edit" src="../assets/svg/edit.svg" alt="edit"/>
                <img @click.stop class="remove" src="../assets/svg/remove.svg" alt="remove" @click="removeRule('routes', index)"/>
              </div>
            </td>
          </tr>
          <tr class="rule-none" v-if="selected.routes != null">
            <td colspan="4">
              <div class="save" @click="saveConfig">
                <img src="../assets/svg/save.svg"/>
              </div>
              <div class="extra">
                <div class="plus">
                  <img src="../assets/svg/plus.svg" @click="addItem('routes', 'host')"/>
                </div>
                <div class="menu">
                  <p class="title">Hostname</p>
                  <div class="edit" v-for="(item, index) in data.routes[selected.routes].host" :key="index">
                    <input class="left" v-model="data.routes[selected.routes].host[index]"/>
                    <img class="remove" src="../assets/svg/remove.svg" alt="remove" @click="removeItem('routes', 'host', index)"/>
                  </div>
                </div>
                <div class="plus">
                  <img src="../assets/svg/plus.svg" @click="addItem('routes','firewall')"/>
                </div>
                <div class="menu">
                  <p class="title">Firewall</p>
                  <div class="edit" v-for="(item, index) in data.routes[selected.routes].firewall" :key="index">
                    <input class="left" v-model="data.routes[selected.routes].firewall[index]"/>
                    <img class="remove" src="../assets/svg/remove.svg" alt="remove" @click="removeItem('routes', 'firewall', index)"/>
                  </div>
                </div>
                <div class="plus"></div>
                <div class="menu">
                  <p class="title">Endpoint</p>
                  <div class="edit">
                    <input v-model="data.routes[selected.routes].endpoint"/>
                  </div>
                </div>
              </div>
            </td>
          </tr>
          <tr class="rule-none">
            <td>
              <p @click="addRule('routes')" class="lefts">+ Add more</p>
            </td>
            <td colspan="2"></td>
            <td>
              <p class="center">{{ data.routes.length }}</p>
            </td>
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
            <th class="short">Actions</th>
          </tr>
        </thead>
        <tbody v-if="data.firewall">
          <tr class="rule" v-for="(item, index) in data.firewall" :key="index" @click="selectRule('firewall', index)">
            <td>
              <select class="select">
                <option>
                  {{ item.name }}
                </option>
              </select>
            </td>
            <td>
              <select class="select">
                <option v-if="!item.block">
                  Allow
                </option>
                <option v-else>
                  Block
                </option>
              </select>
            </td>
            <td>
              <select class="select">
                <option v-for="(item, index) in item.cidr" :key="index" :value="item" :disabled="index !== 0">
                  {{ item }}
                </option>
              </select>
            </td>
            <td>
              <div class="setting">
                <img class="edit" src="../assets/svg/edit.svg" alt="edit"/>
                <img @click.stop class="remove" src="../assets/svg/remove.svg" alt="remove" @click="removeRule('firewall', index)"/>
              </div>
            </td>
          </tr>
          <tr class="rule-none" v-if="selected.firewall != null">
            <td colspan="4">
              <div class="save" @click="saveConfig">
                <img src="../assets/svg/save.svg"/>
              </div>
              <div class="extra">
                <div class="menu">
                  <p class="title">Rulename</p>
                  <div class="edit">
                    <input v-model="data.firewall[selected.firewall].name"/>
                  </div>
                </div>
                <div class="plus"></div>
                <div class="menu">
                  <p class="title">To</p>
                  <select v-model="data.firewall[selected.firewall].block" class="to-select">
                    <option :value="false">Allow</option>
                    <option :value="true">Block</option>
                  </select>
                </div>
                <div class="plus">
                  <img src="../assets/svg/plus.svg" @click="addItem('firewall','cidr')"/>
                </div>
                <div class="menu">
                  <p class="title">CIDR</p>
                  <div class="edit" v-for="(item, index) in data.firewall[selected.firewall].cidr" :key="index">
                    <input class="left" v-model="data.firewall[selected.firewall].cidr[index]"/>
                    <img class="remove" src="../assets/svg/remove.svg" alt="remove" @click="removeItem('firewall', 'cidr', index)"/>
                  </div>
                </div>
                <div class="plus"></div>
              </div>
            </td>
          </tr>
          <tr class="rule-none">
            <td>
              <p @click="addRule('firewall')" class="lefts">+ Add more</p>
            </td>
            <td colspan="2"></td>
            <td>
              <p class="center">{{ data.routes.length }}</p>
            </td>
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