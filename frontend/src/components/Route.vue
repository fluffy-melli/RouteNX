<script setup>
import { useI18n } from 'vue-i18n'
import { ref, onMounted } from 'vue'
import { Chart, LinearScale, CategoryScale, LineElement, PointElement, Title, Tooltip, Legend, LineController } from 'chart.js'
const { locale, t } = useI18n({ useScope: 'global' })
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
const logger = ref([])
const selected = ref({})

const debug = false
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
  const response = await fetch(`${api}/config`)
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
              const hor = time.getHours()
              const min = time.getMinutes()
              const scd = time.getSeconds()
              return hor.toString() + ":" + min.toString() +":"+ scd.toString()
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
  getLogger()
  getTraffic(chart)
  setInterval(() => {
    getLogger()
    getTraffic(chart)
  }, 30 * 1000)
})

async function getTraffic(chart) {
  const response = await fetch(`${api}/traffic`)
  if (response.ok) {
    const trf = await response.json()
    chart.data.labels = trf["Label"].slice(-20)
    chart.data.datasets[0].data = trf["TX"].slice(-20)
    chart.data.datasets[1].data = trf["RX"].slice(-20)
    chart.update()
  } else {
    console.error('Failed to fetch data')
  }
}

async function getLogger() {
  const response = await fetch(`${api}/logger`)
  if (response.ok) {
    const data = await response.json()
    logger.value = data
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
  if (selected.value[name] == index) {
    selected.value = { ...selected.value, [name]: null }
  }
  data.value[name].splice(index, 1)
}

function addRule(name) {
  if (name == 'routes') {
    data.value[name].push(baseProxyRule)
  } else {
    data.value[name].push(baseFireWallRule)
  }
  selected.value[name] = data.value[name].length - 1
}

function removeItem(name, type, index) {
  data.value[name][selected.value[name]][type].splice(index, 1)
}

function addItem(name, type) {
  data.value[name][selected.value[name]][type].push("")
}

</script>

<template>
  <div class="app">
    <div>
      <div>
        <h2 class="title">RouteNX</h2>
      </div>
    </div>
    <div class="line"></div>
    <div class="container">
      <h2 class="title">{{ t('proxyTitle') }}</h2>
      <p class="description">{{ t('proxyDescription') }}</p>
      <table class="styled-table">
        <thead>
          <tr>
            <th class="long-4">{{ t('hostname') }}</th>
            <th class="long-4">{{ t('firewall') }}</th>
            <th class="long-4">{{ t('endpoint') }}</th>
            <th class="short">{{ t('actions') }}</th>
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
            <td colspan="3">
              <div class="save" @click="saveConfig">
                <img src="../assets/svg/save.svg"/>
              </div>
              <div class="extra">
                <div class="plus">
                  <img src="../assets/svg/plus.svg" @click="addItem('routes', 'host')"/>
                </div>
                <div class="menu">
                  <p class="title">{{ t('hostname') }}</p>
                  <div class="edit" v-for="(item, index) in data.routes[selected.routes].host" :key="index">
                    <input class="left" v-model="data.routes[selected.routes].host[index]"/>
                    <img class="remove" src="../assets/svg/remove.svg" alt="remove" @click="removeItem('routes', 'host', index)"/>
                  </div>
                </div>
                <div class="plus">
                  <img src="../assets/svg/plus.svg" @click="addItem('routes','firewall')"/>
                </div>
                <div class="menu">
                  <p class="title">{{ t('firewall') }}</p>
                  <div class="edit" v-for="(item, index) in data.routes[selected.routes].firewall" :key="index">
                    <input class="left" v-model="data.routes[selected.routes].firewall[index]"/>
                    <img class="remove" src="../assets/svg/remove.svg" alt="remove" @click="removeItem('routes', 'firewall', index)"/>
                  </div>
                </div>
                <div class="plus"></div>
                <div class="menu">
                  <p class="title">{{ t('endpoint') }}</p>
                  <div class="edit">
                    <input v-model="data.routes[selected.routes].endpoint"/>
                  </div>
                </div>
              </div>
            </td>
          </tr>
          <tr class="rule-none">
            <td colspan="3">
              <p @click="addRule('routes')" class="lefts">+ {{ t('addmore') }}</p>
            </td>
            <td>
              <p class="left-center">{{ data.routes.length }}</p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="line"></div>
    <div class="container">
      <h2 class="title">{{ t('firewallTitle') }}</h2>
      <p class="description">{{ t('firewallDescription') }}</p>
      <table class="styled-table">
        <thead>
          <tr>
            <th class="long-4">{{ t('rulename') }}</th>
            <th class="long-4">{{ t('packet') }}</th>
            <th class="long-4">CIDR</th>
            <th class="short">{{ t('actions') }}</th>
          </tr>
        </thead>
        <tbody v-if="data.firewall">
          <tr class="rule" v-for="(item, index) in data.firewall" :key="index" @click="selectRule('firewall', index)">
            <td>
              <select class="select" @click.stop>
                <option>
                  {{ item.name }}
                </option>
              </select>
            </td>
            <td>
              <select class="select" @click.stop>
                <option v-if="!item.block">
                  {{ t('allow') }}
                </option>
                <option v-else @click.stop>
                  {{ t('block') }}
                </option>
              </select>
            </td>
            <td>
              <select class="select" @click.stop>
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
            <td colspan="3">
              <div class="save" @click="saveConfig">
                <img src="../assets/svg/save.svg"/>
              </div>
              <div class="extra">
                <div class="menu">
                  <p class="title">{{ t('rulename') }}</p>
                  <div class="edit">
                    <input v-model="data.firewall[selected.firewall].name"/>
                  </div>
                </div>
                <div class="plus"></div>
                <div class="menu">
                  <p class="title">To</p>
                  <select v-model="data.firewall[selected.firewall].block" class="to-select">
                    <option :value="false">
                      {{ t('allow') }}
                    </option>
                    <option :value="true">
                      {{ t('block') }}
                    </option>
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
            <td colspan="3">
              <p @click="addRule('firewall')" class="lefts">+ {{ t('addmore') }}</p>
            </td>
            <td>
              <p class="left-center">{{ data.firewall.length }}</p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="line"></div>
    <div class="container">
      <canvas id="traffic"></canvas>
    </div>
    <div class="line"></div>
    <div class="container">
      <div class="log">
        <pre>{{ t('errorLog') }}</pre>
        <pre>{{ t('time') }} | {{ t('error') }}</pre>
        <pre v-for="(item, index) in logger.error" :key="index">{{ item.time }}: {{ item.error }}</pre>
        <pre></pre>
        <pre>{{ t('firewallBlockLog') }}</pre>
        <pre>{{ t('time') }} | {{ t('host') }} | {{ t('forwardIp') }} | {{ t('originIp') }}</pre>
        <pre v-for="(item, index) in logger.block" :key="index">{{ item.time }}: {{ item.host }} {{ item.forward_ip }} {{ item.origin_ip }}</pre>
      </div>
    </div>
    <div class="line"></div>
    <div class="container">
      <div class="center">
        <select v-model="locale" class="lang">
          <option value="kr">
            KR
          </option>
          <option value="en">
            EN
          </option>
        </select>
      </div>
    </div>
  </div>
</template>