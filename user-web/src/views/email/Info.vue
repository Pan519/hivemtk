<template>
  <div class="email-providers-container">
    <el-table :data="filteredProviders" border stripe style="width: 100%">
      <el-table-column prop="name" ::label="$t('auto.Info.f0213e4c')" width="120" fixed />
      <el-table-column prop="vendor" ::label="$t('auto.Info.5219be7d')" width="120" />
      <el-table-column ::label="$t('auto.Info.9a53b2b8')" width="180">
        <template #default="{ row }">
          <el-tag>{{ row.smtp }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('auto.Info.c76cfefe')" width="120">
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
      <el-table-column :label="$t('auto.Info.d5b8cf0c')" min-width="220">
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
      <el-table-column :label="$t('auto.Info.6f1afa4f')" min-width="200">
        <template #default="{ row }">
          <el-link type="primary" :href="row.registerUrl" target="_blank">
            {{ row.registerUrl }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('auto.Info.2b6bc0f2')" width="120" fixed="right">
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
    vendor: t('auto.Info.f396b465'),
    smtp: 'smtp.163.com',
    ports: [{ value: 465, type: 'ssl' }, { value: 994, type: 'ssl' }],
    authSummary: t('auto.Info.8f9f12a0'),
    authNotes: [
      t('auto.Info.d05c2ee1'),
      t('auto.Info.588ea461'),
      t('auto.Info.a67fb0d6')
    ],
    registerUrl: 'https://mail.163.com'
  },
  {
    id: 2,
    name: '126邮箱',
    vendor: t('auto.Info.f396b465'),
    smtp: 'smtp.126.com',
    ports: [{ value: 465, type: 'ssl' }],
    authSummary: t('auto.Info.54d75cd7'),
    authNotes: [
      t('auto.Info.eddb6326'),
      t('auto.Info.18178a78'),
      t('auto.Info.f6b46a63')
    ],
    registerUrl: 'https://mail.126.com'
  },
  {
    id: 3,
    name: 'Yeah邮箱',
    vendor: t('auto.Info.f396b465'),
    smtp: 'smtp.yeah.net',
    ports: [{ value: 25, type: 'plain' }],
    authSummary: t('auto.Info.160a0976'),
    authNotes: [
      t('auto.Info.1e3678bd'),
      t('auto.Info.778d42d7')
    ],
    registerUrl: 'https://mail.yeah.net'
  },
  {
    id: 4,
    name: 'QQ邮箱',
    vendor: t('auto.Info.19778031'),
    smtp: 'smtp.qq.com',
    ports: [{ value: 465, type: 'ssl' }, { value: 587, type: 'tls' }],
    authSummary: t('auto.Info.f6bc2097'),
    authNotes: [
      t('auto.Info.7b0a2a96'),
      t('auto.Info.63844e43'),
      t('auto.Info.5568fe01')
    ],
    registerUrl: 'https://mail.qq.com'
  },
  {
    id: 5,
    name: t('auto.Info.bf7f333b'),
    vendor: t('auto.Info.62a3c7ee'),
    smtp: 'smtp.sina.com',
    ports: [{ value: 25, type: 'plain' }],
    authSummary: t('auto.Info.81f56ac9'),
    authNotes: [
      t('auto.Info.f9a9f547'),
      t('auto.Info.81b6c3a5')
    ],
    registerUrl: 'https://mail.sina.com.cn'
  },
  {
    id: 6,
    name: '139邮箱',
    vendor: t('auto.Info.a470d2e2'),
    smtp: 'smtp.139.com',
    ports: [{ value: 465, type: 'ssl' }],
    authSummary: t('auto.Info.dada7060'),
    authNotes: [
      t('auto.Info.06ca0bf7'),
      t('auto.Info.bb188f2a'),
      t('auto.Info.e95feba9')
    ],
    registerUrl: 'https://mail.10086.cn'
  },
  {
    id: 7,
    name: 'Gmail',
    vendor: 'Google',
    smtp: 'smtp.gmail.com',
    ports: [{ value: 587, type: 'tls' }],
    authSummary: t('auto.Info.c9b726f6'),
    authNotes: [
      t('auto.Info.79274df4'),
      t('auto.Info.3233b822'),
      t('auto.Info.63c4ca0c')
    ],
    registerUrl: 'https://mail.google.com'
  },
  {
    id: 8,
    name: 'Outlook',
    vendor: 'Microsoft',
    smtp: 'smtp.office365.com',
    ports: [{ value: 587, type: 'tls' }],
    authSummary: t('auto.Info.8867c4ae'),
    authNotes: [
      t('auto.Info.ed753e28'),
      t('auto.Info.21c75016'),
      t('auto.Info.f9f2abd9')
    ],
    registerUrl: 'https://outlook.live.com'
  },
  {
    id: 9,
    name: 'Yahoo',
    vendor: 'Yahoo',
    smtp: 'smtp.mail.yahoo.com',
    ports: [{ value: 465, type: 'ssl' }],
    authSummary: t('auto.Info.f9f2abd9'),
    authNotes: [
      t('auto.Info.def2c9d7'),
      t('auto.Info.f9f2abd9'),
      t('auto.Info.e5595dec')
    ],
    registerUrl: 'https://mail.yahoo.com'
  },
  {
    id: 10,
    name: 'AOL',
    vendor: 'Verizon',
    smtp: 'smtp.aol.com',
    ports: [{ value: 587, type: 'tls' }],
    authSummary: t('auto.Info.81f6fc55'),
    authNotes: [t('auto.Info.e21cbe61'), t('auto.Info.10b4d692')],
    registerUrl: 'https://mail.aol.com'
  },
  {
    id: 11,
    name: 'GMX',
    vendor: 'United Internet',
    smtp: 'smtp.gmx.com',
    ports: [{ value: 465, type: 'ssl' }],
    authSummary: t('auto.Info.d8f5f97f'),
    authNotes: [t('auto.Info.5e85d3c7'), t('auto.Info.97f705e0')],
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
