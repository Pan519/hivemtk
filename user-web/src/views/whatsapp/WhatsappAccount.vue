<template>
  <div class="page">
    <el-card>
      <template #header>
        <div class="card-header">{{ $t('auto.WhatsappAccount.7.3bf9cb') }}</div>
      </template>
      <el-form :inline="true" @submit.prevent>
        <el-form-item :label="$t('auto.WhatsappAccount.1.d7ec2d')">
          <el-input v-model="form.name" :placeholder="$t('auto.WhatsappAccount.2.17327f')" />
        </el-form-item>
        <el-form-item :label="$t('auto.WhatsappAccount.3.2432b5')">
          <el-input v-model="form.remark" :placeholder="$t('auto.WhatsappAccount.3.2432b5')" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="createAccount">{{ $t('auto.WhatsappAccount.8.f0677c') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="accounts" style="width:100%" v-loading="loading">
        <el-table-column prop="name" :label="$t('auto.WhatsappAccount.1.d7ec2d')" />
        <el-table-column prop="remark" :label="$t('auto.WhatsappAccount.3.2432b5')" />
        <el-table-column prop="status" :label="$t('auto.WhatsappAccount.4.3fea7c')" />
        <el-table-column :label="$t('auto.WhatsappAccount.5.2b6bc0')" width="340">
          <template #default="{ row }">
            <el-button size="small" @click="startLogin(row)">{{ $t('auto.WhatsappAccount.9.402d19') }}</el-button>
            <el-button size="small" @click="refreshStatus(row)">{{ $t('auto.WhatsappAccount.4.3fea7c') }}</el-button>
            <el-button size="small" type="primary" @click="openBindingDialog(row)">{{ $t('auto.WhatsappAccount.10.c9e87d') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <AgentBindingDialog
      v-model="bindingDialogVisible"
      channel-type="whatsapp"
      :account-id="bindingAccountId"
      :account-label="bindingAccountLabel"
      :account-enabled="true"
    />

    <el-dialog v-model="qrDialog" :title="$t('auto.WhatsappAccount.6.d85b03')" width="360px">
      <div v-if="qrCode" class="qr-box">
        <img :src="qrImageSrc" alt="QR" width="300" height="300" />
      </div>
      <div v-else>{{ $t('auto.WhatsappAccount.11.866d40') }}</div>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api/whatsapp'
import { ElMessage } from 'element-plus'
import { toList } from '@/utils/list'
import { computed } from 'vue'

const { t } = useI18n()

const accounts = ref([])
const loading = ref(false)
const form = ref({ name: '', remark: '' })
const qrDialog = ref(false)
const qrCode = ref('')
const qrImageSrc = computed(() => qrCode.value ? `https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=${encodeURIComponent(qrCode.value)}` : '')

onMounted(() => { loadAccounts() })

function loadAccounts() {
  loading.value = true
  api.listAccounts().then(res => {
    accounts.value = toList(res)
  }).finally(() => loading.value = false)
}

function createAccount() {
  if (!form.value.name) { ElMessage.error(i18n.global.t(t('auto.WhatsappAccount.12.06e2f8'))); return }
  api.createAccount(form.value).then(() => {
    ElMessage.success(i18n.global.t(t('auto.WhatsappAccount.13.04a691')))
    form.value = { name: '', remark: '' }
    loadAccounts()
  })
}

function startLogin(row) {
  api.startLogin(row.id).then(res => {
    qrCode.value = res.qr;
    qrDialog.value = true
  }).catch((e) => {
    ElMessage.error(e?.message || t('auto.WhatsappAccount.14.7515ac'));
  })
}

function refreshStatus(row) {
  api.loginStatus(row.id).then(res => {
    if (res.logged_in) {
      ElMessage.success(i18n.global.t(t('auto.WhatsappAccount.15.71fa3b')))
      loadAccounts()
    } else {
      qrCode.value = res.qr || ''
      if (qrCode.value) qrDialog.value = true
      else ElMessage.info(i18n.global.t(t('auto.WhatsappAccount.16.e6499a')))
    }
  }).catch((e) => {
    ElMessage.error(e?.message || t('auto.WhatsappAccount.17.540141'))
  })
}

const bindingDialogVisible = ref(false);
const bindingAccountId = ref(0)
const bindingAccountLabel = ref('')
function openBindingDialog(row) {
  bindingAccountId.value = row.id
  bindingAccountLabel.value = row.name || row.remark || `#${row.id}`
  bindingDialogVisible.value = true
}
</script>

<style scoped>
.page { padding: 20px; }
.card-header { font-weight: 600; }
.qr-box { display:flex; justify-content:center; }
</style>
