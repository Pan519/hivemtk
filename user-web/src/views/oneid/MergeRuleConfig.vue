<template>
  <div class="merge-rule-config">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.MergeRuleConfig.2f83b354') }}</span>
          <el-button type="primary" @click="addRule">{{ $t('auto.MergeRuleConfig.1d0c7fb6') }}</el-button>
        </div>
      </template>
      <el-alert type="info" :closable="false" style="margin-bottom: 16px">
        <p>{{ $t('auto.MergeRuleConfig.eed00855') }}</p>
        <p>{{ $t('auto.MergeRuleConfig.5e4b62f0') }}</p>
      </el-alert>

      <el-table :data="rules" v-loading="loading">
        <el-table-column prop="name" :label="$t('auto.MergeRuleConfig.87080256')" width="180" />
        <el-table-column prop="type" :label="$t('auto.MergeRuleConfig.226b0912')" width="120">
          <template #default="{ row }">
            <el-tag :type="row.type === 'deterministic' ? 'success' : 'warning'">
              {{ row.type === 'deterministic' ? '确定性' : '概率性' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="fields" :label="$t('auto.MergeRuleConfig.ed0bc856')" min-width="200">
          <template #default="{ row }">
            <el-tag v-for="f in row.fields" :key="f" size="small" style="margin-right: 4px">
              {{ fieldLabel(f) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="threshold" :label="$t('auto.MergeRuleConfig.bd493e1b')" width="140">
          <template #default="{ row }">
            <el-progress :percentage="row.threshold" :stroke-width="6" />
          </template>
        </el-table-column>
        <el-table-column prop="priority" :label="$t('auto.MergeRuleConfig.ee8ecb9e')" width="80" />
        <el-table-column :label="$t('auto.MergeRuleConfig.7854b52a')" width="80">
          <template #default="{ row }">
            <el-switch v-model="row.enabled" />
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.MergeRuleConfig.2b6bc0f2')" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editRule(row)">{{ $t('auto.MergeRuleConfig.95b351c8') }}</el-button>
            <el-button link type="danger" @click="deleteRule(row)">{{ $t('auto.MergeRuleConfig.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-dialog v-model="dialogVisible" :title="editing.id ? '编辑规则' : '新增规则'" width="640px">
      <el-form :model="editing" label-width="100px">
        <el-form-item :label="$t('auto.MergeRuleConfig.87080256')">
          <el-input v-model="editing.name" />
        </el-form-item>
        <el-form-item :label="$t('auto.MergeRuleConfig.226b0912')">
          <el-radio-group v-model="editing.type">
            <el-radio label="deterministic">{{ $t('auto.MergeRuleConfig.2d08d2f0') }}</el-radio>
            <el-radio label="probabilistic">{{ $t('auto.MergeRuleConfig.e16dbb5e') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('auto.MergeRuleConfig.ed0bc856')">
          <el-checkbox-group v-model="editing.fields">
            <el-checkbox label="phone">{{ $t('auto.MergeRuleConfig.8098e2b4') }}</el-checkbox>
            <el-checkbox label="email">{{ $t('auto.MergeRuleConfig.3bc5e602') }}</el-checkbox>
            <el-checkbox label="wechat_open_id">{{ $t('auto.MergeRuleConfig.147f0337') }}</el-checkbox>
            <el-checkbox label="douyin_id">{{ $t('auto.MergeRuleConfig.0807cf97') }}</el-checkbox>
            <el-checkbox label="nickname">{{ $t('auto.MergeRuleConfig.23eb0e60') }}</el-checkbox>
            <el-checkbox label="region">{{ $t('auto.MergeRuleConfig.2560b304') }}</el-checkbox>
            <el-checkbox label="device_fingerprint">{{ $t('auto.MergeRuleConfig.039fbcca') }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item :label="$t('auto.MergeRuleConfig.bd493e1b')">
          <el-slider v-model="editing.threshold" :min="50" :max="100" :step="5" show-input />
        </el-form-item>
        <el-form-item :label="$t('auto.MergeRuleConfig.ee8ecb9e')">
          <el-input-number v-model="editing.priority" :min="0" :max="100" />
        </el-form-item>
        <el-form-item :label="$t('auto.MergeRuleConfig.c500cfab')">
          <el-radio-group v-model="editing.action">
            <el-radio label="auto">{{ $t('auto.MergeRuleConfig.b3e2ba86') }}</el-radio>
            <el-radio label="review">{{ $t('auto.MergeRuleConfig.ea316c8c') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.MergeRuleConfig.625fb26b') }}</el-button>
        <el-button type="primary" @click="saveRule">{{ $t('auto.MergeRuleConfig.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>

    
    <el-card v-if="preview" class="preview-card">
      <template #header>
        <span>{{ $t('auto.MergeRuleConfig.97c52012') }}</span>
      </template>
      <el-statistic :value="preview.candidateCount" :title="$t('auto.MergeRuleConfig.5fcd3e06')" />
      <el-table :data="preview.samples" size="small" max-height="240">
        <el-table-column prop="from" label="From" />
        <el-table-column prop="to" label="To" />
        <el-table-column prop="score" :label="$t('auto.MergeRuleConfig.678cd563')" width="120">
          <template #default="{ row }">
            <el-progress :percentage="row.score" :stroke-width="6" />
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n';
import { ElMessage, ElMessageBox } from 'element-plus'
import { http } from '@/utils/request'

const FIELD_LABELS = {
  phone: i18n.global.t('auto.MergeRuleConfig.8098e2b4'),
  email: i18n.global.t('auto.MergeRuleConfig.3bc5e602'),
  wechat_open_id: i18n.global.t('auto.MergeRuleConfig.cfbf6f4c'),
  douyin_id: i18n.global.t('auto.MergeRuleConfig.277a984f'),
  nickname: i18n.global.t('auto.MergeRuleConfig.23eb0e60'),
  region: i18n.global.t('auto.MergeRuleConfig.2560b304'),
  device_fingerprint: i18n.global.t('auto.MergeRuleConfig.039fbcca')
}
const fieldLabel = (k) => FIELD_LABELS[k] || k

const { t } = useI18n()

const loading = ref(false)
const rules = ref([])
const dialogVisible = ref(false)
const editing = reactive({ id: null, name: '', type: 'deterministic', fields: [], threshold: 95, priority: 50, action: 'review', enabled: true })
const preview = ref(null)

async function load() {
  loading.value = true
  try {
    const res = await http.get('/api/oneid/merge-rules')
    const list = Array.isArray(res)
      ? res
      : [...(res?.built_in || []), ...(res?.custom || [])];
    rules.value = list.map(r => ({
      ...r,
      type: r.type || 'deterministic',
      fields: Array.isArray(r.fields) ? r.fields : (r.field ? [r.field] : []),
      threshold: typeof r.threshold === 'number' ? r.threshold : 100
    }))
  } finally {
    loading.value = false
  }
}

async function saveRule() {
  if (!editing.name) {
    ElMessage.warning(t('auto.MergeRuleConfig.411cac4a'))
    return
  }
  if (editing.fields.length === 0) {
    ElMessage.warning(t('auto.MergeRuleConfig.9397ace4'))
    return
  }
  await http.post('/api/oneid/merge-rules', editing)
  ElMessage.success(t('auto.MergeRuleConfig.f8dfedcd'))
  dialogVisible.value = false
  await load()
  await runPreview()
}

function addRule() {
  Object.assign(editing, { id: null, name: '', type: 'deterministic', fields: [], threshold: 95, priority: 50, action: 'review', enabled: true })
  dialogVisible.value = true
}

function editRule(row) {
  Object.assign(editing, row)
  dialogVisible.value = true
}

async function deleteRule(row) {
  try {
    await ElMessageBox.confirm(`确认删除规则「${row.name}」？`, t('auto.MergeRuleConfig.50eaf94d'), { type: 'warning' })
  } catch (e) {
    if (e === 'cancel' || e === 'close') return
    throw e
  }
  await http.delete(`/api/oneid/merge-rules/${row.id}`)
  ElMessage.success(t('auto.MergeRuleConfig.5cc23262'))
  await load()
}

async function runPreview() {
  try {
    preview.value = await http.post('/api/oneid/merge-rules/preview', rules.value)
  } catch (_) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}

onMounted(async () => {
  await load()
  await runPreview()
})
</script>

<style scoped>
.merge-rule-config { padding: 16px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.preview-card { margin-top: 16px; }
</style>
