<template>
  <div class="llm-routing-page">
    <el-card class="header-card">
      <div class="header-content">
        <h2>{{ $t('auto.List.fdc6011c') }}</h2>
        <p class="subtitle">{{ $t('auto.List.a402fea9') }}</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showModelDialog()">
          <el-icon><Plus /></el-icon>
          新增模型
        </el-button>
        <el-button @click="$router.push('/llmRouting/cost')">{{ $t('auto.List.99da45bc') }}</el-button>
        <el-button @click="refreshAll">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </el-card>

    <el-tabs v-model="activeTab" class="content-tabs">
      
      <el-tab-pane :label="$t('auto.List.851f3dff')" name="models">
        <el-table :data="models" v-loading="loading.models" stripe>
          <template #empty>
            <el-empty :description="$t('auto.List.285a2e77')" />
          </template>
          <el-table-column prop="name" :label="$t('auto.List.920fe38e')" min-width="140" />
          <el-table-column prop="vendor" :label="$t('auto.List.70aefc6a')" width="110">
            <template #default="{ row }">
              <el-tag size="small" :type="row.vendor === 'local' ? 'success' : 'info'">
                {{ vendorLabel(row.vendor) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="model" :label="$t('auto.List.f97c1388')" min-width="160" show-overflow-tooltip />
          <el-table-column :label="$t('auto.List.3fea7ca7')" width="90">
            <template #default="{ row }">
              <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
                {{ row.enabled ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="max_rpm" :label="$t('auto.List.e8718247')" width="110" align="center" />
          <el-table-column prop="cost_per_1k" :label="$t('auto.List.94b3d898')" width="120" align="center">
            <template #default="{ row }">
              <span :class="row.cost_per_1k > 0 ? 'cost-paid' : 'cost-free'">
                {{ row.cost_per_1k > 0 ? row.cost_per_1k.toFixed(4) : '免费' }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="base_url" :label="$t('auto.List.f682f85e')" min-width="200" show-overflow-tooltip />
          <el-table-column label="API Key" width="90">
            <template #default="{ row }">
              <el-tag v-if="row.api_key_set" type="success" size="small">{{ $t('auto.List.44e60789') }}</el-tag>
              <el-tag v-else type="info" size="small">{{ $t('auto.List.fe2d26a2') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.bd177f3a')" width="80">
            <template #default="{ row }">
              <el-tag v-if="row.no_fc" type="warning" size="small">NoFC</el-tag>
              <el-tag v-else type="success" size="small">{{ $t('auto.List.262b0b07') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.2b6bc0f2')" width="260" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="testModel(row)">{{ $t('auto.List.db06c78d') }}</el-button>
              <el-button link type="primary" @click="toggleStatus(row)">
                {{ row.enabled ? '禁用' : '启用' }}
              </el-button>
              <el-button link type="primary" @click="showModelDialog(row)">{{ $t('auto.List.95b351c8') }}</el-button>
              <el-button link type="danger" @click="deleteModel(row)">{{ $t('auto.List.2f4aaddd') }}</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.List.603ec630')" name="routing">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.78862aed') }}</span>
              <div>
                <el-button size="small" @click="loadRouting">
                  <el-icon><Refresh /></el-icon> 刷新
                </el-button>
                <el-button type="primary" size="small" @click="addRoutingRow">
                  <el-icon><Plus /></el-icon> 新增映射
                </el-button>
              </div>
            </div>
          </template>
          <el-table :data="sceneRouting" v-loading="loading.routing" stripe>
            <template #empty>
              <el-empty :description="$t('auto.List.de4dec26')" />
            </template>
            <el-table-column prop="scenario" :label="$t('auto.List.c931653c')" min-width="160" />
            <el-table-column prop="provider" :label="$t('auto.List.85b70c4a')" min-width="140" />
            <el-table-column :label="$t('auto.List.6a8a0069')" min-width="200">
              <template #default="{ row }">
                <el-tag v-for="fb in (row.fallbacks || [])" :key="fb" size="small" type="info" style="margin-right: 4px">
                  {{ fb }}
                </el-tag>
                <span v-if="!row.fallbacks?.length" class="text-muted">-</span>
              </template>
            </el-table-column>
            <el-table-column prop="version" :label="$t('auto.List.fe2df04a')" width="80" align="center" />
            <el-table-column prop="weight" :label="$t('auto.List.e4455fd1')" width="110" align="center">
              <template #default="{ row }">
                <el-progress
                  v-if="row.weight > 0 && row.weight < 100"
                  :percentage="Number(row.weight)"
                  :stroke-width="10"
                  :color="canaryColors"
                />
                <el-tag v-else-if="row.weight === 100" type="warning" size="small">{{ $t('auto.List.0d3c7a08') }}</el-tag>
                <span v-else class="text-muted">-</span>
              </template>
            </el-table-column>
            <el-table-column prop="min_quality" :label="$t('auto.List.0f851629')" width="100" align="center" />
            <el-table-column prop="max_latency" :label="$t('auto.List.19510ce7')" width="120" align="center" />
            <el-table-column :label="$t('auto.List.2b6bc0f2')" width="160" fixed="right">
              <template #default="{ row, $index }">
                <el-button link type="primary" @click="editRoutingRow(row, $index)">{{ $t('auto.List.95b351c8') }}</el-button>
                <el-button link type="danger" @click="removeRoutingRow($index)">{{ $t('auto.List.2f4aaddd') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.List.96264b6b')" name="fallback">
        <el-card v-loading="loading.fallback">
          <template #header><span>{{ $t('auto.List.bbbabb7b') }}</span></template>
          <el-table :data="fallbackRows" stripe>
            <template #empty>
              <el-empty :description="$t('auto.List.612531ae')" />
            </template>
            <el-table-column prop="scenario" :label="$t('auto.List.c931653c')" min-width="160" />
            <el-table-column :label="$t('auto.List.2b935d09')" min-width="140">
              <template #default="{ row }">
                <el-tag type="primary" size="small">{{ row.provider }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.0208fce1')" min-width="300">
              <template #default="{ row }">
                <template v-if="row.fallbacks?.length">
                  <el-tag v-for="(fb, i) in row.fallbacks" :key="fb" size="small" type="info" style="margin-right: 4px">
                    {{ i + 1 }}. {{ fb }}
                  </el-tag>
                </template>
                <span v-else class="text-muted">{{ $t('auto.List.69a24bad') }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="version" :label="$t('auto.List.fe2df04a')" width="80" align="center" />
          </el-table>
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.List.380e60e2')" name="cost">
        <el-row :gutter="20" class="stat-row">
          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-label">{{ $t('auto.List.7f6f35c4') }}</div>
              <div class="stat-value">{{ costStats.total_calls || 0 }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-label">{{ $t('auto.List.808bcc57') }}</div>
              <div class="stat-value" style="color: #4F46E5">{{ formatNumber(costStats.total_tokens) }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-label">{{ $t('auto.List.d6ee0789') }}</div>
              <div class="stat-value" style="color: #EF4444">{{ (costStats.monthly_cost || 0).toFixed(4) }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-label">{{ $t('auto.List.5691eda9') }}</div>
              <div class="stat-value" style="color: #F59E0B">{{ costStats.enabled_models || 0 }} / {{ costStats.active_models || 0 }}</div>
            </el-card>
          </el-col>
        </el-row>
        <el-card style="margin-top: 20px">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.f2547cf6') }}</span>
              <el-radio-group v-model="costWindow" size="small" @change="loadCost">
                <el-radio-button label="today">{{ $t('auto.List.296304f1') }}</el-radio-button>
                <el-radio-button label="week">{{ $t('auto.List.d765efc4') }}</el-radio-button>
                <el-radio-button label="month">{{ $t('auto.List.c059bed7') }}</el-radio-button>
                <el-radio-button label="all">{{ $t('auto.List.a8b0c204') }}</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <el-table :data="costStats.by_provider || []" v-loading="loading.cost" stripe>
            <template #empty><el-empty :description="$t('auto.List.1e091b12')" /></template>
            <el-table-column prop="provider" label="Provider" min-width="140" />
            <el-table-column prop="call_count" :label="$t('auto.List.00d5bdf1')" width="120" align="center" />
            <el-table-column prop="total_tokens" :label="$t('auto.List.8952bc8d')" width="140" align="center" />
            <el-table-column prop="total_cost" :label="$t('auto.List.1f9507e1')" width="120" align="center">
              <template #default="{ row }">{{ (row.total_cost || 0).toFixed(4) }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.c909e6eb')" width="200" align="center">
              <template #default="{ row }">
                <el-progress
                  :percentage="calcRatio(row.total_cost)"
                  :stroke-width="10"
                />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
        <el-card style="margin-top: 20px">
          <template #header><span>{{ $t('auto.List.8fb479ae') }}</span></template>
          <el-table :data="costStats.by_scenario || []" stripe>
            <template #empty><el-empty :description="$t('auto.List.6ed55112')" /></template>
            <el-table-column prop="scenario" :label="$t('auto.List.c931653c')" min-width="160" />
            <el-table-column prop="call_count" :label="$t('auto.List.00d5bdf1')" width="120" align="center" />
            <el-table-column prop="total_tokens" :label="$t('auto.List.8952bc8d')" width="140" align="center" />
            <el-table-column prop="total_cost" :label="$t('auto.List.1f9507e1')" width="120" align="center">
              <template #default="{ row }">{{ (row.total_cost || 0).toFixed(4) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.List.39c6c930')" name="audit">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.88fac83f') }}</span>
              <el-button size="small" @click="loadAudit">
                <el-icon><Refresh /></el-icon> 刷新
              </el-button>
            </div>
          </template>
          <el-table :data="auditList" v-loading="loading.audit" stripe>
            <template #empty><el-empty :description="$t('auto.List.79c05840')" /></template>
            <el-table-column :label="$t('auto.List.19fcb9eb')" min-width="180">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
            <el-table-column prop="scenario" :label="$t('auto.List.c931653c')" min-width="140" />
            <el-table-column prop="version" :label="$t('auto.List.fe2df04a')" width="80" align="center" />
            <el-table-column :label="$t('auto.List.1dd554cd')" min-width="280">
              <template #default="{ row }">
                <span v-if="row.action === 'update_strategy'">
                  <el-tag size="small" type="info">{{ row.prev_provider || '无' }}</el-tag>
                  →
                  <el-tag size="small" type="success">{{ row.new_provider }}</el-tag>
                </span>
                <el-tag v-else-if="row.action === 'create_model'" type="primary" size="small">{{ $t('auto.List.1cbddef0') }}</el-tag>
                <el-tag v-else-if="row.action === 'delete_model'" type="danger" size="small">{{ $t('auto.List.bb7a2ea8') }}</el-tag>
                <el-tag v-else-if="row.action === 'update_model'" type="warning" size="small">{{ $t('auto.List.bb1abfa5') }}</el-tag>
                <el-tag v-else size="small">{{ getRoutingActionLabel(row.action) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="operator" :label="$t('auto.List.f9ac4b2a')" width="120" />
            <el-table-column prop="trace_id" label="Trace ID" min-width="180" show-overflow-tooltip />
          </el-table>
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.List.c6b03576')" name="modeltype">
        <el-row :gutter="20" class="stat-row">
          <el-col :span="12">
            <el-card v-loading="loading.modelType">
              <template #header>
                <div class="card-header">
                  <span>{{ $t('auto.List.ce20e7a4') }}</span>
                  <el-radio-group v-model="modelTypeWindow" size="small" @change="loadModelType">
                    <el-radio-button label="today">{{ $t('auto.List.296304f1') }}</el-radio-button>
                    <el-radio-button label="week">{{ $t('auto.List.d765efc4') }}</el-radio-button>
                    <el-radio-button label="month">{{ $t('auto.List.c059bed7') }}</el-radio-button>
                    <el-radio-button label="all">{{ $t('auto.List.a8b0c204') }}</el-radio-button>
                  </el-radio-group>
                </div>
              </template>
              <el-row :gutter="12">
                <el-col :span="8">
                  <div class="metric-box local">
                    <div class="metric-label">{{ $t('auto.List.9b91fcc5') }}</div>
                    <div class="metric-value">{{ getModelTypeCount('local') }}</div>
                    <div class="metric-sub">¥{{ getModelTypeCost('local').toFixed(4) }}</div>
                  </div>
                </el-col>
                <el-col :span="8">
                  <div class="metric-box cloud">
                    <div class="metric-label">{{ $t('auto.List.c29b95b1') }}</div>
                    <div class="metric-value">{{ getModelTypeCount('cloud') }}</div>
                    <div class="metric-sub">¥{{ getModelTypeCost('cloud').toFixed(4) }}</div>
                  </div>
                </el-col>
                <el-col :span="8">
                  <div class="metric-box">
                    <div class="metric-label">{{ $t('auto.List.3c95425a') }}</div>
                    <div class="metric-value">{{ getModelTypeCount('missing') }}</div>
                    <div class="metric-sub">{{ getMissingRatio() }}%</div>
                  </div>
                </el-col>
              </el-row>
              <el-divider />
              <div class="self-sufficiency">
                <div class="metric-label">本地自给率（{{ formatPercent(getSelfSufficiency()) }}%）</div>
                <el-progress
                  :percentage="Number(formatPercent(getSelfSufficiency()))"
                  :stroke-width="14"
                  :color="['#67c23a', '#e6a23c', '#f56c6c']"
                />
              </div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card v-loading="loading.health">
              <template #header>
                <div class="card-header">
                  <span>{{ $t('auto.List.37b20722') }}</span>
                  <el-button size="small" @click="loadHealth">
                    <el-icon><Refresh /></el-icon> 刷新
                  </el-button>
                </div>
              </template>
              <el-table :data="llmHealth.providers || []" stripe size="small">
                <template #empty><el-empty :description="$t('auto.List.b856d0cf')" /></template>
                <el-table-column prop="name" label="Provider" min-width="120" />
                <el-table-column prop="vendor" :label="$t('auto.List.70aefc6a')" width="100" />
                <el-table-column :label="$t('auto.List.3fea7ca7')" width="100" align="center">
                  <template #default="{ row }">
                    <el-tag v-if="row.healthy" type="success" size="small">{{ $t('auto.List.fbc026c9') }}</el-tag>
                    <el-tag v-else-if="row.circuit_open" type="danger" size="small">{{ $t('auto.List.3a0c7850') }}</el-tag>
                    <el-tag v-else type="warning" size="small">{{ $t('auto.List.2981d4c8') }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column :label="$t('auto.List.7030ff64')" width="80" align="center">
                  <template #default="{ row }">{{ row.error_count || 0 }}</template>
                </el-table-column>
                <el-table-column prop="avg_latency" :label="$t('auto.List.32da342a')" width="100" align="center" />
                <el-table-column :label="$t('auto.List.c7efa1d6')" min-width="160" show-overflow-tooltip>
                  <template #default="{ row }">{{ row.last_error || '-' }}</template>
                </el-table-column>
              </el-table>
            </el-card>
          </el-col>
        </el-row>
        <el-card style="margin-top: 20px" v-loading="loading.egress">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.182ec9d6') }}</span>
              <el-button size="small" @click="loadEgressAlerts">
                <el-icon><Refresh /></el-icon> 刷新
              </el-button>
            </div>
          </template>
          <el-alert
            v-if="egressAlerts.length === 0"
            type="success"
            :closable="false"
            :title="$t('auto.List.6ba8c159')"
            :description="$t('auto.List.c4ca1a0c')"
            show-icon
          />
          <el-table v-else :data="egressAlerts" stripe size="small">
            <el-table-column prop="scenario" :label="$t('auto.List.c931653c')" min-width="120" />
            <el-table-column prop="provider" label="Provider" min-width="120" />
            <el-table-column prop="model" label="Model" min-width="140" />
            <el-table-column prop="base_url" :label="$t('auto.List.aa5c2464')" min-width="240" show-overflow-tooltip />
            <el-table-column prop="vendor" :label="$t('auto.List.70aefc6a')" width="100" />
            <el-table-column :label="$t('auto.List.f672c2c9')" width="100" align="center">
              <template #default="{ row }">{{ row.call_count || 0 }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.ac59bc30')" min-width="180">
              <template #default="{ row }">{{ formatTime(row.first_seen) }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.8f10f0ef')" min-width="180">
              <template #default="{ row }">{{ formatTime(row.last_seen) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    
    <el-dialog v-model="modelDialogVisible" :title="modelDialogTitle" width="640px">
      <el-form :model="modelForm" :rules="modelFormRules" ref="modelFormRef" label-width="100px">
        <el-form-item :label="$t('auto.List.5d551c2c')" prop="name">
          <el-input v-model="modelForm.name" :placeholder="$t('auto.List.0f89de58')" :disabled="!!modelForm._isEdit" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.f97c1388')" prop="model">
          <el-input v-model="modelForm.model" :placeholder="$t('auto.List.a1ba1b31')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.f682f85e')" prop="base_url">
          <el-input v-model="modelForm.base_url" :placeholder="$t('auto.List.a38eb0cf')" />
        </el-form-item>
        <el-form-item label="API Key" prop="api_key">
          <el-input v-model="modelForm.api_key" type="password" show-password :placeholder="$t('auto.List.5a23035c')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3ea3a5c7')">
          <el-input-number v-model="modelForm.quality_score" :min="0" :max="1" :step="0.01" :precision="2" />
          <span class="form-tip">{{ $t('auto.List.b824660f') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.List.e8718247')">
          <el-input-number v-model="modelForm.max_rpm" :min="0" :step="10" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.94b3d898')">
          <el-input-number v-model="modelForm.cost_per_1k" :min="0" :precision="4" :step="0.001" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.7854b52a')">
          <el-switch v-model="modelForm.enabled" />
        </el-form-item>
        <el-form-item label="NoFC">
          <el-switch v-model="modelForm.no_fc" />
          <span class="form-tip">{{ $t('auto.List.faf48eb0') }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="modelDialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" @click="submitModel">{{ $t('auto.List.38cf16f2') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="routingDialogVisible" :title="routingDialogTitle" width="600px">
      <el-form :model="routingForm" :rules="routingFormRules" ref="routingFormRef" label-width="120px">
        <el-form-item :label="$t('auto.List.c931653c')" prop="scenario">
          <el-select v-model="routingForm.scenario" :placeholder="$t('auto.List.354a40ac')" filterable allow-create style="width: 100%">
            <el-option v-for="s in scenarioPresets" :key="s" :label="s" :value="s" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.85b70c4a')" prop="provider">
          <el-select v-model="routingForm.provider" :placeholder="$t('auto.List.f2d3731b')" style="width: 100%">
            <el-option v-for="m in models" :key="m.name" :label="`${m.name} (${m.vendor})`" :value="m.name" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.6a8a0069')">
          <el-select v-model="routingForm.fallbacks" multiple style="width: 100%" :placeholder="$t('auto.List.02a77f5d')">
            <el-option v-for="m in models" :key="m.name" :label="m.name" :value="m.name" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.e4455fd1')">
          <el-input-number v-model="routingForm.weight" :min="0" :max="100" :step="5" />
          <span class="form-tip">{{ $t('auto.List.d5e4ce6b') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.List.0f851629')">
          <el-input-number v-model="routingForm.min_quality" :min="0" :max="1" :step="0.01" :precision="2" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.19510ce7')">
          <el-input-number v-model="routingForm.max_latency" :min="0" :step="1000" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="routingDialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" @click="submitRouting">{{ $t('auto.List.38cf16f2') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import { LlmRoutingApi } from '@/api/llmRouting.js'

const { t } = useI18n()

const activeTab = ref('models')
const loading = reactive({ models: false, routing: false, fallback: false, cost: false, audit: false, modelType: false, health: false, egress: false })

const models = ref([]);
const sceneRouting = ref([])
const costStats = ref({})
const auditList = ref([])
const costWindow = ref('month')

const modelTypeWindow = ref('month');
const modelTypeStats = ref({ by_provider_type: [], self_sufficiency: 0, total_calls: 0 })
const llmHealth = ref({ status: 'unknown', providers: [], healthy_providers: 0, total_providers: 0 })
const egressAlerts = ref([])

const fallbackRows = computed(() => sceneRouting.value.map(r => ({
  scenario: r.scenario,
  provider: r.provider,
  fallbacks: r.fallbacks || [],
  version: r.version,
})));

const scenarioPresets = [
  'intent_recognize', 'sop_reply', 'objection', 'friendly_chat',
  'long_summary', 'high_quality', 'low_cost',
]

const canaryColors = [
  { color: '#909399', percentage: 30 },
  { color: '#F59E0B', percentage: 60 },
  { color: '#EF4444', percentage: 100 },
]

const modelDialogVisible = ref(false);
const modelDialogTitle = ref(t('auto.List.1cbddef0'))
const modelFormRef = ref()
const emptyModelForm = () => ({
  name: '', model: '', base_url: '', api_key: '',
  quality_score: 0.95, max_rpm: 60, cost_per_1k: 0,
  enabled: true, no_fc: false, _isEdit: false,
})
const modelForm = ref(emptyModelForm())
const modelFormRules = {
  name: [{ required: true, message: t('auto.List.97b6dccd'), trigger: 'blur' }],
  model: [{ required: true, message: t('auto.List.4b6ff7d4'), trigger: 'blur' }],
  base_url: [{ required: true, message: t('auto.List.8007271c'), trigger: 'blur' }],
}

const routingDialogVisible = ref(false);
const routingDialogTitle = ref(t('auto.List.374e0375'))
const routingFormRef = ref()
const emptyRoutingForm = () => ({
  scenario: '',
  provider: '',
  fallbacks: [],
  weight: 0,
  min_quality: 0.8,
  max_latency: 90000,
})
const routingForm = ref(emptyRoutingForm())
const routingEditIndex = ref(-1)
const routingFormRules = {
  scenario: [{ required: true, message: t('auto.List.75f77abd'), trigger: 'change' }],
  provider: [{ required: true, message: t('auto.List.e347fa73'), trigger: 'change' }],
}

const loadModels = async () => {
  loading.models = true
  try {
    const res = await LlmRoutingApi.getModelList()
    models.value = Array.isArray(res) ? res : (res?.list || res?.items || res?.data || [])
  } catch (e) {
    models.value = []
  } finally {
    loading.models = false
  }
};

const loadRouting = async () => {
  loading.routing = true
  try {
    const res = await LlmRoutingApi.getSceneRouting()
    sceneRouting.value = Array.isArray(res) ? res : (res?.list || res?.items || res?.data || [])
  } catch (e) {
    sceneRouting.value = []
  } finally {
    loading.routing = false
  }
}

const loadFallback = async () => {
  loading.fallback = true
  try {
    if (sceneRouting.value.length === 0) {
      await loadRouting()
    }
  } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ } finally {
    loading.fallback = false
  }
}

const loadCost = async () => {
  loading.cost = true
  try {
    const res = await LlmRoutingApi.getCostStats(costWindow.value)
    costStats.value = res || {}
  } catch (e) {
    costStats.value = {}
  } finally {
    loading.cost = false
  }
}

const loadAudit = async () => {
  loading.audit = true
  try {
    const res = await LlmRoutingApi.getAuditHistory(null, 50)
    auditList.value = Array.isArray(res) ? res : (res?.list || res?.items || res?.data || [])
  } catch (e) {
    auditList.value = []
  } finally {
    loading.audit = false
  }
}

const loadModelType = async () => {
  loading.modelType = true
  try {
    const res = await LlmRoutingApi.getModelTypeStats(modelTypeWindow.value)
    modelTypeStats.value = res || { by_provider_type: [], self_sufficiency: 0, total_calls: 0 }
  } catch (e) {
    modelTypeStats.value = { by_provider_type: [], self_sufficiency: 0, total_calls: 0 }
  } finally {
    loading.modelType = false
  }
};

const loadHealth = async () => {
  loading.health = true
  try {
    const res = await LlmRoutingApi.getHealth()
    llmHealth.value = res || { status: 'unknown', providers: [], healthy_providers: 0, total_providers: 0 }
  } catch (e) {
    llmHealth.value = { status: 'unknown', providers: [], healthy_providers: 0, total_providers: 0 }
  } finally {
    loading.health = false
  }
}

const loadEgressAlerts = async () => {
  loading.egress = true
  try {
    const res = await LlmRoutingApi.getEgressAlerts()
    egressAlerts.value = Array.isArray(res) ? res : (res?.list || res?.items || res?.data || [])
  } catch (e) {
    egressAlerts.value = []
  } finally {
    loading.egress = false
  }
}

const findModelTypeRow = (modelType) => {
  const arr = modelTypeStats.value.by_provider_type || []
  return arr.find(r => r.model_type === modelType) || {}
};

const getModelTypeCount = (modelType) => {
  return Number(findModelTypeRow(modelType).call_count || 0)
}

const getModelTypeCost = (modelType) => {
  return Number(findModelTypeRow(modelType).total_cost || 0)
}

const getSelfSufficiency = () => {
  const raw = Number(modelTypeStats.value.self_sufficiency || 0);
  return raw <= 1 ? raw * 100 : raw
}

const getMissingRatio = () => {
  const total = Number(modelTypeStats.value.total_calls || 0)
  if (total === 0) return '0.00'
  const missing = getModelTypeCount('missing')
  return ((missing / total) * 100).toFixed(2)
}

const formatPercent = (v) => {
  const n = Number(v || 0)
  return Math.min(100, Math.max(0, n)).toFixed(2)
}

const refreshAll = () => {
  loadModels()
  loadRouting()
  loadCost()
  loadAudit()
  loadModelType()
  loadHealth()
  loadEgressAlerts()
}

const showModelDialog = (row) => {
  if (row) {
    modelForm.value = {
      ...emptyModelForm(),
      name: row.name,
      model: row.model,
      base_url: row.base_url,
      api_key: '',
      quality_score: row.quality_score,
      max_rpm: row.max_rpm,
      cost_per_1k: row.cost_per_1k,
      enabled: row.enabled,
      no_fc: row.no_fc,
      _isEdit: true,
    }
    modelDialogTitle.value = t('auto.List.c85b90dd')
  } else {
    modelForm.value = emptyModelForm()
    modelDialogTitle.value = t('auto.List.1cbddef0')
  }
  modelDialogVisible.value = true
};

const submitModel = async () => {
  if (!modelFormRef.value) return
  await modelFormRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const payload = { ...modelForm.value }
      delete payload._isEdit
      if (modelForm.value._isEdit) {
        await LlmRoutingApi.updateModel(modelForm.value.name, payload);
        ElMessage.success(t('auto.List.55aa6366'))
      } else {
        await LlmRoutingApi.createModel(payload)
        ElMessage.success(t('auto.List.a5bfd70d'))
      }
      modelDialogVisible.value = false
      loadModels()
    } catch (e) {
      ElMessage.error(e?.message || t('auto.List.6de920b4'))
    }
  })
}

const deleteModel = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除模型 "${row.name}" 吗？`, t('auto.List.e83a256e'), { type: 'warning' })
    await LlmRoutingApi.deleteModel(row.name)
    ElMessage.success(t('auto.List.0007d170'))
    loadModels()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(t('auto.List.acf0664a'))
  }
}

const toggleStatus = async (row) => {
  const next = !row.enabled
  try {
    await LlmRoutingApi.updateModel(row.name, { enabled: next });
    ElMessage.success(t('auto.List.7d6e0e1a'))
    loadModels()
  } catch (e) {
    ElMessage.error(t('auto.List.effbdce5') + (e?.message || ''))
  }
}

const testModel = async (row) => {
  try {
    const result = await LlmRoutingApi.testModel(row.name, { prompt: t('auto.List.7f370af4'), timeout_seconds: 30 })
    if (result?.success) {
      ElMessage.success(`模型 "${row.name}" 连通性正常（延迟 ${result.latency_ms}ms）`)
    } else {
      ElMessage.warning(`模型 "${row.name}" 调用失败：${result?.error || t('auto.List.974e7484')}`)
    }
  } catch (e) {
    ElMessage.error(t('auto.List.4b3cb9fe') + (e?.message || ''))
  }
}

const addRoutingRow = () => {
  routingForm.value = emptyRoutingForm()
  routingEditIndex.value = -1
  routingDialogTitle.value = t('auto.List.374e0375')
  routingDialogVisible.value = true
};

const editRoutingRow = (row, index) => {
  routingForm.value = {
    scenario: row.scenario,
    provider: row.provider,
    fallbacks: [...(row.fallbacks || [])],
    weight: row.weight || 0,
    min_quality: row.min_quality || 0.8,
    max_latency: row.max_latency || 90000,
  }
  routingEditIndex.value = index
  routingDialogTitle.value = t('auto.List.cb14eb12')
  routingDialogVisible.value = true
}

const removeRoutingRow = async (index) => {
  try {
    await ElMessageBox.confirm(t('auto.List.b9adc1be'), t('auto.List.e83a256e'), { type: 'warning' })
    const removed = sceneRouting.value[index]
    const next = sceneRouting.value.filter((_, i) => i !== index)
    sceneRouting.value = next
    await LlmRoutingApi.saveSceneRouting(next);
    ElMessage.success(`已删除 ${removed?.scenario || ''}`)
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(t('auto.List.acf0664a'))
  }
}

const submitRouting = async () => {
  if (!routingFormRef.value) return
  await routingFormRef.value.validate(async (valid) => {
    if (!valid) return
    const next = [...sceneRouting.value]
    if (routingEditIndex.value >= 0) {
      next[routingEditIndex.value] = { ...routingForm.value }
    } else {
      next.push({ ...routingForm.value })
    }
    try {
      await LlmRoutingApi.saveSceneRouting(next)
      ElMessage.success(t('auto.List.3b108349'))
      routingDialogVisible.value = false
      loadRouting()
    } catch (e) {
      ElMessage.error(t('auto.List.78d235cc') + (e?.message || ''))
    }
  })
}

const vendorLabel = (v) => ({
  local: t('auto.List.8e7eca1f'), deepseek: 'DeepSeek', qwen: t('auto.List.4680d368'),
  openai: 'OpenAI', zhipu: t('auto.List.7f567eea'), moonshot: t('auto.List.e562c560'), other: t('auto.List.5a0afc9d'),
}[v] || v || '-');

const formatNumber = (n) => {
  if (!n) return 0
  return Number(n).toLocaleString()
}

const formatTime = (s) => {
  if (!s) return '-'
  return new Date(s).toLocaleString('zh-CN')
}

const calcRatio = (cost) => {
  if (!cost || !costStats.value.monthly_cost) return 0
  const total = costStats.value.monthly_cost;
  if (total <= 0) return 0
  return Math.min(100, Math.round((cost / total) * 100))
}

onMounted(() => {
  refreshAll()
})
</script>

<style scoped lang="scss">
.llm-routing-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .header-content h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
  .header-actions { display: flex; gap: 10px; }
}
.content-tabs { background: #fff; padding: 16px; border-radius: 4px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.stat-row { margin-bottom: 20px; }
.stat-card {
  text-align: center;
  .stat-label { color: #909399; font-size: 14px; margin-bottom: 10px; }
  .stat-value { font-size: 28px; font-weight: bold; }
}
.form-tip { color: #909399; font-size: 12px; margin-left: 12px; }
.text-muted { color: #c0c4cc; }
.cost-free { color: #67C23A; }
.cost-paid { color: #F56C6C; }

// v3.7.0：分类统计 / 自给率 metric-box
.metric-box {
  padding: 16px 8px;
  border-radius: 8px;
  background: #f7f8fa;
  text-align: center;
  margin-bottom: 8px;
  .metric-label { color: #909399; font-size: 13px; }
  .metric-value { font-size: 28px; font-weight: 700; color: #303133; margin: 4px 0; }
  .metric-sub { color: #909399; font-size: 12px; }
  &.local { background: #f0f9eb; }
  &.cloud { background: #fdf6ec; }
}
.self-sufficiency { padding: 0 4px; }
</style>
