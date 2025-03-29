<script setup>
import { ref, onMounted } from 'vue'

const data = ref([])

onMounted(async () => {
  const response = await fetch('http://localhost:3000/route')
  if (response.ok) {
    data.value = await response.json()
  } else {
    console.error('Failed to fetch data')
  }
})

function getSettingClass(name, list, index) {
  if (index === 0 && index !== list.length - 1) {
    return `${name}-start`;
  } else if (index !== 0 && index === list.length - 1) {
    return `${name}-end`;
  } else {
    return `${name}-all`;
  }
}
</script>

<template>
  <div class="app">
    <div class="container">
      <h2>Reverse Proxy</h2>
      <div v-for="(item, index) in data.routes" :key="index" :class="getSettingClass('route', data.routes, index)">
        <div class="info">
          <select>
            <option v-for="(item, index) in item.host" :key="index" :value="item" :disabled="index !== 0">
              {{ item }}
            </option>
          </select>
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
          <select>
            <option>
              {{ item.endpoint }}
            </option>
          </select>
        </div>
        <div :class="getSettingClass('setting', data.routes, index)">
          <img class="remove" src="../assets/svg/remove.svg" alt="remove"/>
          <img class="setting" src="../assets/svg/setting.svg" alt="setting"/>
        </div>
      </div>
    </div>
    <div class="container">
      <h2>Firewall</h2>
      <div v-for="(item, index) in data.firewall" :key="index" :class="getSettingClass('route', data.firewall, index)">
        <div class="info">
          <select>
            <option>
              {{ item.name }}
            </option>
          </select>
          <select>
            <option v-for="(item, index) in item.cidr" :key="index" :value="item" :disabled="index !== 0">
              {{ item }}
            </option>
          </select>
          <select>
            <option v-if="!item.block">
              Allow
            </option>
            <option v-else>
              Block
            </option>
          </select>
        </div>
        <div :class="getSettingClass('setting', data.firewall, index)">
          <img class="remove" src="../assets/svg/remove.svg" alt="remove"/>
          <img class="setting" src="../assets/svg/setting.svg" alt="setting"/>
        </div>
      </div>
    </div>
  </div>
</template>