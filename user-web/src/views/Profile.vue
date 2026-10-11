<template>
  <div class="profile-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('个人资料') }}</h2>
          <p class="subtitle">{{ $t('查看与维护当前登录账号的基础信息') }}</p>
        </div>
        <div class="header-actions">
          <el-button :loading="loading" @click="loadProfile">
            <el-icon><Refresh /></el-icon>
            {{ $t('刷新') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <el-row :gutter="20" v-loading="loading">
      <el-col :xs="24" :sm="10" :md="8">
        <el-card class="avatar-card" shadow="never">
          <div class="avatar-wrap">
            <el-avatar :size="96" class="avatar">{{ initial }}</el-avatar>
            <div class="role-tag">
              <el-tag :type="roleTagType" effect="dark" size="small">
                {{ roleLabel }}
              </el-tag>
            </div>
            <h3 class="username">{{ form.username || '-' }}</h3>
            <p class="user-id">ID: {{ form.id || '-' }}</p>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="14" :md="16">
        <el-card shadow="never">
          <template #header>
            <span>{{ $t('基础信息') }}</span>
          </template>
          <el-form
            :model="form"
            :rules="rules"
            ref="formRef"
            label-width="100px"
            label-position="right"
          >
            <el-form-item :label="$t('auto.Profile.1.bb2cdf')" prop="username">
              <el-input v-model="form.username" disabled />
            </el-form-item>
            <el-form-item :label="$t('auto.Profile.2.7eacb4')" prop="real_name">
              <el-input v-model="form.real_name" :placeholder="$t('auto.Profile.3.3348b7')" maxlength="32" />
            </el-form-item>
            <el-form-item :label="$t('auto.Profile.4.3bc5e6')" prop="email">
              <el-input v-model="form.email" :placeholder="$t('auto.Profile.5.dbf6d0')" maxlength="64" />
            </el-form-item>
            <el-form-item :label="$t('auto.Profile.6.8098e2')" prop="phone">
              <el-input v-model="form.phone" :placeholder="$t('auto.Profile.7.6e4f4b')" maxlength="20" />
            </el-form-item>
            <el-form-item :label="$t('auto.Profile.8.464f3d')">
              <el-tag :type="roleTagType">{{ roleLabel }}</el-tag>
              <span class="readonly-hint">{{ $t('auto.Profile.18.5b6eb1') }}</span>
            </el-form-item>
            <el-form-item :label="$t('auto.Profile.9.3fea7c')">
              <el-tag v-if="form.status === 1" type="success">{{ $t('auto.Profile.19.7854b5') }}</el-tag>
              <el-tag v-else type="info">{{ $t('auto.Profile.20.69b0f6') }}</el-tag>
            </el-form-item>
            <el-form-item :label="$t('auto.Profile.10.eca37c')">
              <span class="readonly-text">{{ form.created_at || '-' }}</span>
            </el-form-item>
            <el-form-item :label="$t('auto.Profile.11.19fb21')">
              <span class="readonly-text">{{ form.last_login_at || '尚未记录' }}</span>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="onSave">
                <el-icon><Check /></el-icon>
                保存修改
              </el-button>
              <el-button @click="loadProfile">{{ $t('auto.Profile.21.4b9c32') }}</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card class="password-card" shadow="never">
          <template #header>
            <span>{{ $t('auto.Profile.22.7fc88a') }}</span>
          </template>
          <el-alert
            type="info"
            :closable="false"
            show-icon
            style="margin-bottom: 16px"
          >
            修改密码后需要使用新密码重新登录。
          </el-alert>
          <el-form
            ref="pwdFormRef"
            :model="pwdForm"
            :rules="pwdRules"
            label-width="100px"
            label-position="right"
          >
            <el-form-item :label="$t('auto.Profile.12.a806ce')" prop="oldPassword">
              <el-input
                v-model="pwdForm.oldPassword"
                type="password"
                show-password
                :placeholder="$t('auto.Profile.13.cd043c')"
              />
            </el-form-item>
            <el-form-item :label="$t('auto.Profile.14.bf7da0')" prop="newPassword">
              <el-input
                v-model="pwdForm.newPassword"
                type="password"
                show-password
                :placeholder="$t('auto.Profile.15.204d3e')"
              />
            </el-form-item>
            <el-form-item :label="$t('auto.Profile.16.670d63')" prop="confirmPassword">
              <el-input
                v-model="pwdForm.confirmPassword"
                type="password"
                show-password
                :placeholder="$t('auto.Profile.17.ec7b73')"
                @keyup.enter="onChangePassword"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="warning" :loading="pwdSaving" @click="onChangePassword">
                <el-icon><Lock /></el-icon>
                修改密码
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Check, Lock } from '@element-plus/icons-vue'
import { http } from '@/utils/request'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
void Check; void Lock

const { t } = useI18n()

const loading = ref(false)
const saving = ref(false)
const pwdSaving = ref(false)

const formRef = ref(null)
const pwdFormRef = ref(null)

const form = reactive({
  id: '',
  username: '',
  real_name: '',
  email: '',
  phone: '',
  role: '',
  status: 1,
  created_at: '',
  last_login_at: ''
})

const pwdForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const initial = computed(() => (form.username || '?').charAt(0).toUpperCase())

const roleLabel = computed(() => {
  const r = form.role || userStore.role || ''
  if (r === 'admin') return t('auto.Profile.23.1dd8e0')
  if (r === 'manager') return t('auto.Profile.24.c0477f')
  if (r === 'sales') return t('auto.Profile.25.4ae98a')
  if (r === 'viewer') return t('auto.Profile.26.6c5190')
  return r || t('auto.Profile.27.fe2d26')
})

const roleTagType = computed(() => {
  const r = form.role || userStore.role
  if (r === 'admin') return 'danger'
  if (r === 'manager') return 'warning'
  if (r === 'sales') return 'success'
  return 'info'
})

const rules = {
  real_name: [{ max: 32, message: i18n.global.t(t('auto.Profile.28.9864ba')), trigger: 'blur' }],
  email: [{
    validator: (_rule, value, cb) => {
      if (!value) return cb()
      if (!/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(value)) {
        return cb(new Error(t('auto.Profile.29.75262e')))
      }
      cb()
    },
    trigger: 'blur'
  }],
  phone: [{
    validator: (_rule, value, cb) => {
      if (!value) return cb()
      if (!/^1[3-9]\d{9}$/.test(value)) {
        return cb(new Error(t('auto.Profile.30.19c896')))
      }
      cb()
    },
    trigger: 'blur'
  }]
}

const pwdRules = {
  oldPassword: [{ required: true, message: i18n.global.t(t('auto.Profile.13.cd043c')), trigger: 'blur' }],
  newPassword: [
    { required: true, message: i18n.global.t(t('auto.Profile.31.abdd7e')), trigger: 'blur' },
    { min: 8, max: 64, message: i18n.global.t(t('auto.Profile.32.73ca94')), trigger: 'blur' },
    {
      validator: (_rule, value, cb) => {
        if (!value) return cb()
        if (!/[a-z]/.test(value)) return cb(new Error(t('auto.Profile.33.799739')))
        if (!/[A-Z]/.test(value)) return cb(new Error(t('auto.Profile.34.bb02c8')))
        if (!/\d/.test(value)) return cb(new Error(t('auto.Profile.35.a65c9e')))
        cb()
      },
      trigger: 'blur'
    }
  ],
  confirmPassword: [
    { required: true, message: i18n.global.t(t('auto.Profile.36.d98a14')), trigger: 'blur' },
    {
      validator: (_rule, value, cb) => {
        if (value !== pwdForm.newPassword) {
          return cb(new Error(t('auto.Profile.37.6046d6')))
        }
        cb()
      },
      trigger: 'blur'
    }
  ]
}

const loadProfile = async () => {
  loading.value = true
  try {
    const resp = await http.get('/api/auth/current-user')
    const data = resp || {}
    Object.assign(form, {
      id: data.id || '',
      username: data.username || '',
      real_name: data.real_name || data.realName || '',
      email: data.email || '',
      phone: data.phone || data.mobile || '',
      role: data.role || '',
      status: data.status ?? 1,
      created_at: data.created_at || data.createdAt || '',
      last_login_at: data.last_login_at || data.lastLoginAt || ''
    })
    userStore.setUserInfo({
      id: data.id,
      username: data.username,
      email: data.email,
      role: data.role
    });
  } catch (err) {
    ElMessage.error(t('auto.Profile.38.19a682') + (err?.message || err))
  } finally {
    loading.value = false
  }
};

const onSave = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    saving.value = true
    await http.put(`/api/users/${form.id}`, {
      real_name: form.real_name,
      email: form.email,
      phone: form.phone
    })
    ElMessage.success(i18n.global.t(t('auto.Profile.39.3b1083')))
  } catch (err) {
    if (err?.message) {
      ElMessage.error(t('auto.Profile.40.78d235') + err.message)
    }
  } finally {
    saving.value = false
  }
};

const onChangePassword = async () => {
  if (!pwdFormRef.value) return
  try {
    await pwdFormRef.value.validate()
    try {
      await ElMessageBox.confirm(
        t('auto.Profile.41.0abe7e'),
        t('auto.Profile.42.146713'),
        { type: 'warning', confirmButtonText: t('auto.Profile.43.49e56c'), cancelButtonText: t('auto.Profile.44.625fb2') }
      )
    } catch {
      return;
    }
    pwdSaving.value = true
    await http.post('/api/auth/change-password', {
      old_password: pwdForm.oldPassword,
      new_password: pwdForm.newPassword
    })
    ElMessage.success(i18n.global.t(t('auto.Profile.45.4e400f')))
    setTimeout(() => {
      userStore.logout()
      window.location.href = '/#/login'
    }, 1200)
  } catch (err) {
    if (err?.message) {
      ElMessage.error(t('auto.Profile.46.c26a1c') + err.message)
    }
  } finally {
    pwdSaving.value = false
  }
};

onMounted(() => {
  loadProfile()
})
</script>

<style scoped>
.profile-page {
  padding: 0;
}
.header-card {
  margin-bottom: 16px;
}
.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-content h2 {
  margin: 0 0 4px;
  font-size: 20px;
}
.subtitle {
  color: #909399;
  font-size: 13px;
  margin: 0;
}
.avatar-card {
  text-align: center;
}
.avatar-wrap {
  padding: 20px 0;
}
.avatar {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff;
  font-size: 36px;
  font-weight: 600;
}
.role-tag {
  margin: 12px 0 8px;
}
.username {
  margin: 8px 0 4px;
  font-size: 18px;
  color: #303133;
}
.user-id {
  margin: 0;
  color: #909399;
  font-size: 12px;
}
.password-card {
  margin-top: 16px;
}
.readonly-hint {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}
.readonly-text {
  color: #606266;
  font-size: 13px;
}
</style>
