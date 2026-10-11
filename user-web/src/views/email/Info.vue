<template>
  <div class="email-providers-container">
    <el-table :data="filteredProviders" border stripe style="width: 100%">
      <el-table-column prop="name" ::label="$t('auto.Info.1.f0213e')" width="120" fixed />
      <el-table-column prop="vendor" ::label="$t('auto.Info.2.5219be')" width="120" />
      <el-table-column ::label="$t('auto.Info.3.9a53b2')" width="180">
        <template #default="{ row }">
          <el-tag>{{ row.smtp }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('auto.Info.4.c76cfe')" width="120">
        <template #default="{ row }">
          <div class="port-badges">
            <el-tag
              v-for="port in row.ports"
              :key="port.value"
              :type="getPortType(port.type)"
            >
              {{ port.value }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('auto.Info.5.d5b8cf')" min-width="220">
        <template #default="{ row }">
          <div class="auth-notes">
            <el-popover placement="top" width="300" trigger="hover">
              <template #reference>
                <el-icon :size="18" style="margin-right: 8px">
                  <InfoFilled />
                </el-icon>
              </template>
              <div class="popover-content">
                <p v-for="(note, idx) in row.authNotes" :key="idx">
                  {{ note }}
                </p>
              </div>
            </el-popover>
            <span>{{ row.authSummary }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('auto.Info.6.6f1afa')" min-width="200">
        <template #default="{ row }">
          <el-link type="primary" :href="row.registerUrl" target="_blank">
            {{ row.registerUrl }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('auto.Info.7.2b6bc0')" width="120" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="goRegister(row)">
            前往注册
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Search, InfoFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const { t } = useI18n()

const searchText = ref('')

const emailProviders = ref([
  {
    id: 1,
    name: '163邮箱',
    vendor: t('auto.Info.8.f396b4'),
    smtp: 'smtp.163.com',
    ports: [{ value: 465, type: 'ssl' }, { value: 994, type: 'ssl' }],
    authSummary: t('auto.Info.9.8f9f12'),
    authNotes: [
      t('auto.Info.10.d05c2e'),
      t('auto.Info.11.588ea4'),
      t('auto.Info.12.a67fb0')
    ],
    registerUrl: 'https://mail.163.com'
  },
  {
    id: 2,
    name: '126邮箱',
    vendor: t('auto.Info.8.f396b4'),
    smtp: 'smtp.126.com',
    ports: [{ value: 465, type: 'ssl' }],
    authSummary: t('auto.Info.13.54d75c'),
    authNotes: [
      t('auto.Info.14.eddb63'),
      t('auto.Info.15.18178a'),
      t('auto.Info.16.f6b46a')
    ],
    registerUrl: 'https://mail.126.com'
  },
  {
    id: 3,
    name: 'Yeah邮箱',
    vendor: t('auto.Info.8.f396b4'),
    smtp: 'smtp.yeah.net',
    ports: [{ value: 25, type: 'plain' }],
    authSummary: t('auto.Info.17.160a09'),
    authNotes: [
      t('auto.Info.18.1e3678'),
      t('auto.Info.19.778d42')
    ],
    registerUrl: 'https://mail.yeah.net'
  },
  {
    id: 4,
    name: 'QQ邮箱',
    vendor: t('auto.Info.20.197780'),
    smtp: 'smtp.qq.com',
    ports: [{ value: 465, type: 'ssl' }, { value: 587, type: 'tls' }],
    authSummary: t('auto.Info.21.f6bc20'),
    authNotes: [
      t('auto.Info.22.7b0a2a'),
      t('auto.Info.23.63844e'),
      t('auto.Info.24.5568fe')
    ],
    registerUrl: 'https://mail.qq.com'
  },
  {
    id: 5,
    name: t('auto.Info.25.bf7f33'),
    vendor: t('auto.Info.26.62a3c7'),
    smtp: 'smtp.sina.com',
    ports: [{ value: 25, type: 'plain' }],
    authSummary: t('auto.Info.27.81f56a'),
    authNotes: [
      t('auto.Info.28.f9a9f5'),
      t('auto.Info.29.81b6c3')
    ],
    registerUrl: 'https://mail.sina.com.cn'
  },
  {
    id: 6,
    name: '139邮箱',
    vendor: t('auto.Info.30.a470d2'),
    smtp: 'smtp.139.com',
    ports: [{ value: 465, type: 'ssl' }],
    authSummary: t('auto.Info.31.dada70'),
    authNotes: [
      t('auto.Info.32.06ca0b'),
      t('auto.Info.33.bb188f'),
      t('auto.Info.34.e95feb')
    ],
    registerUrl: 'https://mail.10086.cn'
  },
  {
    id: 7,
    name: 'Gmail',
    vendor: 'Google',
    smtp: 'smtp.gmail.com',
    ports: [{ value: 587, type: 'tls' }],
    authSummary: t('auto.Info.35.c9b726'),
    authNotes: [
      t('auto.Info.36.79274d'),
      t('auto.Info.37.3233b8'),
      t('auto.Info.38.63c4ca')
    ],
    registerUrl: 'https://mail.google.com'
  },
  {
    id: 8,
    name: 'Outlook',
    vendor: 'Microsoft',
    smtp: 'smtp.office365.com',
    ports: [{ value: 587, type: 'tls' }],
    authSummary: t('auto.Info.39.8867c4'),
    authNotes: [
      t('auto.Info.40.ed753e'),
      t('auto.Info.41.21c750'),
      t('auto.Info.42.f9f2ab')
    ],
    registerUrl: 'https://outlook.live.com'
  },
  {
    id: 9,
    name: 'Yahoo',
    vendor: 'Yahoo',
    smtp: 'smtp.mail.yahoo.com',
    ports: [{ value: 465, type: 'ssl' }],
    authSummary: t('auto.Info.42.f9f2ab'),
    authNotes: [
      t('auto.Info.43.def2c9'),
      t('auto.Info.42.f9f2ab'),
      t('auto.Info.44.e5595d')
    ],
    registerUrl: 'https://mail.yahoo.com'
  },
  {
    id: 10,
    name: 'AOL',
    vendor: 'Verizon',
    smtp: 'smtp.aol.com',
    ports: [{ value: 587, type: 'tls' }],
    authSummary: t('auto.Info.45.81f6fc'),
    authNotes: [t('auto.Info.46.e21cbe'), t('auto.Info.47.10b4d6')],
    registerUrl: 'https://mail.aol.com'
  },
  {
    id: 11,
    name: 'GMX',
    vendor: 'United Internet',
    smtp: 'smtp.gmx.com',
    ports: [{ value: 465, type: 'ssl' }],
    authSummary: t('auto.Info.48.d8f5f9'),
    authNotes: [t('auto.Info.49.5e85d3'), t('auto.Info.50.97f705')],
    registerUrl: 'https://www.gmx.com'
  }
]);

const filteredProviders = computed(() => {
  if (!searchText.value) return emailProviders.value
  const query = searchText.value.toLowerCase()
  return emailProviders.value.filter(
    provider =>
      provider.name.toLowerCase().includes(query) &&
      provider.vendor.toLowerCase().includes(query) &&
      provider.smtp.toLowerCase().includes(query)
  )
});

const getPortType = type => {
  const types = { ssl: 'success', tls: 'warning', plain: 'info' }
  return types[type] || ''
};

const goRegister = provider => {    
  window.open(provider.registerUrl, '_blank')
};
</script>

<style scoped>
.email-providers-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 10px;
}

.search-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.port-badges .el-tag {
  margin-right: 5px;
  margin-bottom: 5px;
}

.auth-notes {
  display: flex;
  align-items: center;
}

.notes-section {
  margin-top: 25px;
}

.notes-section ul {
  margin: 10px 0 0 20px;
}

.popover-content p {
  margin: 5px 0;
  line-height: 1.5;
}
</style>
