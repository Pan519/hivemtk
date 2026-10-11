<template>
  <div class="segment-builder">
    <el-row :gutter="16">
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.Builder.b0fd0de2') }}</span>
              <el-button-group>
                <el-button @click="addNode('and')">{{ $t('auto.Builder.41cd0bf4') }}</el-button>
                <el-button @click="addNode('or')">{{ $t('auto.Builder.86143c1a') }}</el-button>
                <el-button @click="addNode('condition')">{{ $t('auto.Builder.52fa34e1') }}</el-button>
              </el-button-group>
            </div>
          </template>
          <div class="rule-tree">
            <RuleNode v-for="(node, i) in rules" :key="i" :node="node" :level="0" @update="(n) => rules[i] = n" @remove="removeNode(i)" />
          </div>
          <div class="sql-preview">
            <h4>{{ $t('auto.Builder.64816d1d') }}</h4>
            <pre>{{ sqlPreview }}</pre>
            <el-statistic :value="sqlCount" :title="$t('auto.Builder.63db3897')" />
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header><span>{{ $t('auto.Builder.4696724e') }}</span></template>
          <el-form :model="trigger" label-width="80px">
            <el-form-item :label="$t('auto.Builder.159dbc2f')">
              <el-radio-group v-model="trigger.type">
                <el-radio label="manual">{{ $t('auto.Builder.2a3e7f5c') }}</el-radio>
                <el-radio label="schedule">{{ $t('auto.Builder.72ebfe28') }}</el-radio>
                <el-radio label="event">{{ $t('auto.Builder.79ff3e87') }}</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item v-if="trigger.type === 'schedule'" label="Cron">
              <el-input v-model="trigger.cron" placeholder="0 0 9 * * ?" />
            </el-form-item>
            <el-form-item v-if="trigger.type === 'event'" :label="$t('auto.Builder.10b2761d')">
              <el-select v-model="trigger.event">
                <el-option :label="$t('auto.Builder.0f1aa9c0')" value="customer.created" />
                <el-option :label="$t('auto.Builder.c1a80c76')" value="order.completed" />
                <el-option :label="$t('auto.Builder.beaf9a03')" value="behavior.high_value" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('auto.Builder.d081f249')">
              <el-select v-model="trigger.sink" multiple>
                <el-option :label="$t('auto.Builder.58dfae26')" value="reach" />
                <el-option :label="$t('auto.Builder.33ae22a7')" value="cs" />
                <el-option :label="$t('auto.Builder.ea43925b')" value="bulk" />
              </el-select>
            </el-form-item>
          </el-form>
          <el-button type="primary" style="width: 100%" @click="save">{{ $t('auto.Builder.797f7ef0') }}</el-button>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, h } from 'vue'
import { useI18n } from 'vue-i18n';
import { ElMessage } from 'element-plus'
import { http } from '@/utils/request'
// RuleNode 此前只在模板里用、从未 import 也从未定义（git HEAD 同样如此），
// 导致 Vue 每次渲染都报 Failed to resolve component: RuleNode，且整棵规则树 UI 根本不显示。
// 数据结构由本文件的 addNode / compileNode 唯一定义，组件按该契约实现。
import RuleNode from './RuleNode.vue'

const { t } = useI18n()

const rules = ref([
  {
    type: 'and',
    children: [
      { type: 'condition', field: 'last_order_at', op: '>', value: '30d' },
      { type: 'condition', field: 'lifetime_value', op: '>=', value: 1000 }
    ]
  }
])

const trigger = reactive({
  type: 'manual',
  cron: '',
  event: '',
  sink: ['reach']
})

const sqlPreview = computed(() => buildSQL(rules.value))
const sqlCount = ref(0)

function addNode(type) {
  if (type === 'and' || type === 'or') {
    rules.value.push({ type, children: [] })
  } else {
    rules.value.push({ type: 'condition', field: '', op: '=', value: '' })
  }
}

function removeNode(idx) {
  rules.value.splice(idx, 1)
}

function buildSQL(nodes) {
  const parts = nodes.map(compileNode).filter(Boolean);
  return `SELECT id, name FROM customers WHERE ${parts.join(' AND ')}`
}

function compileNode(node) {
  if (node.type === 'and') return `(${node.children?.map(compileNode).join(' AND ') || '1=1'})`
  if (node.type === 'or') return `(${node.children?.map(compileNode).join(' OR ') || '1=1'})`
  if (node.type === 'condition') {
    if (!node.field) return null
    return `${node.field} ${node.op} '${node.value}'`
  }
  return null
}

async function save() {
  // 后端 SegmentSaveRequest.Trigger 是 string；trigger 对象序列化为 JSON 字符串存放
  await http.post('/api/user-segments', {
    name: `分群 ${new Date().toISOString()}`,
    rules: rules.value,
    trigger: JSON.stringify(trigger)
  })
  ElMessage.success(t('auto.Builder.e82641a7'))
}
</script>

<style scoped>
.segment-builder { padding: 16px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.rule-tree { min-height: 240px; }
.sql-preview { margin-top: 16px; padding: 12px; background: #F8FAFC; border-radius: 6px; }
.sql-preview pre { font-size: 12px; color: #334155; white-space: pre-wrap; }
</style>
