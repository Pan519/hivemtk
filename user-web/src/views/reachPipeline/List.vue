<template>
  <div class="reach-pipeline-page">
    <el-tabs v-model="activeTab" @tab-change="handleTabChange" class="rp-tabs">
      
      <el-tab-pane ::label="$t('auto.List.1.5aaab9')" name="pipelines">
        
        <el-row :gutter="16" class="stats-row">
          <el-col :xs="12" :sm="12" :md="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-icon stat-icon-blue"><el-icon><Connection /></el-icon></div>
              <div class="stat-body">
                <div class="stat-title">{{ $t('auto.List.76.6f0393') }}</div>
                <div class="stat-value">{{ stats.total ?? 0 }}</div>
                <div class="stat-sub">活跃 {{ stats.active ?? 0 }} · 暂停 {{ stats.paused ?? 0 }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :xs="12" :sm="12" :md="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-icon stat-icon-green"><el-icon><VideoPlay /></el-icon></div>
              <div class="stat-body">
                <div class="stat-title">{{ $t('运行中') }}</div>
                <div class="stat-value">{{ stats.active ?? 0 }}</div>
                <div class="stat-sub">{{ $t('auto.List.77.790074') }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :xs="12" :sm="12" :md="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-icon stat-icon-orange"><el-icon><VideoPause /></el-icon></div>
              <div class="stat-body">
                <div class="stat-title">{{ $t('已暂停') }}</div>
                <div class="stat-value">{{ stats.paused ?? 0 }}</div>
                <div class="stat-sub">已归档 {{ archivedCount }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :xs="12" :sm="12" :md="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-icon stat-icon-purple"><el-icon><FolderOpened /></el-icon></div>
              <div class="stat-body">
                <div class="stat-title">{{ $t('已归档') }}</div>
                <div class="stat-value">{{ archivedCount }}</div>
                <div class="stat-sub">{{ $t('总数减活跃与暂停') }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :xs="12" :sm="12" :md="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-icon stat-icon-cyan"><el-icon><List /></el-icon></div>
              <div class="stat-body">
                <div class="stat-title">{{ $t('任务总数') }}</div>
                <div class="stat-value">{{ stats.jobs ?? 0 }}</div>
                <div class="stat-sub">待执行 {{ stats.pending ?? 0 }} · 运行中 {{ stats.running ?? 0 }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :xs="12" :sm="12" :md="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-icon stat-icon-green2"><el-icon><CircleCheck /></el-icon></div>
              <div class="stat-body">
                <div class="stat-title">{{ $t('成功数') }}</div>
                <div class="stat-value">{{ stats.success ?? 0 }}</div>
                <div class="stat-sub">成功率 {{ overallSuccessRate }}%</div>
              </div>
            </el-card>
          </el-col>
          <el-col :xs="12" :sm="12" :md="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-icon stat-icon-red"><el-icon><CircleClose /></el-icon></div>
              <div class="stat-body">
                <div class="stat-title">{{ $t('失败数') }}</div>
                <div class="stat-value">{{ stats.failed ?? 0 }}</div>
                <div class="stat-sub">失败率 {{ overallFailureRate }}%</div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <el-card>
          <div class="toolbar">
            <div class="toolbar-left">
              <el-input
                v-model="pipeFilter.keyword"
                ::placeholder="$t('auto.List.2.c6cd09')"
                clearable
                style="width: 180px"
                @keyup.enter="onPipeFilterChange"
                @clear="onPipeFilterChange"
              />
              <el-select
                v-model="pipeFilter.channel"
                ::placeholder="$t('auto.List.3.349bd7')"
                clearable
                style="width: 140px; margin-left: 8px"
                @change="onPipeFilterChange"
              >
                <el-option v-for="c in CHANNEL_OPTIONS" :key="c.value" :label="c.label" :value="c.value" />
              </el-select>
              <el-select
                v-model="pipeFilter.status"
                ::placeholder="$t('auto.List.4.1b8578')"
                clearable
                style="width: 140px; margin-left: 8px"
                @change="onPipeFilterChange"
              >
                <el-option ::label="$t('auto.List.5.61269f')" value="active" />
                <el-option ::label="$t('auto.List.6.aea836')" value="paused" />
                <el-option ::label="$t('auto.List.7.8e327a')" value="archived" />
              </el-select>
              <el-button style="margin-left: 8px" @click="resetPipeFilter">{{ $t('重置') }}</el-button>
              <el-button type="warning" plain @click="resetRateLimit">{{ $t('重置限流') }}</el-button>
            </div>
            <el-button style="margin-left: 8px" @click="$router.push('/reachPipeline/editor')">
              <el-icon><Share /></el-icon>&nbsp;可视化编辑器
            </el-button>
            <el-button type="primary" @click="openCreatePipe">
              <el-icon><Plus /></el-icon>&nbsp;新建 Pipeline
            </el-button>
          </div>

          <el-table :data="filteredPipelines" v-loading="pipeLoading" stripe>
            <el-table-column prop="name" ::label="$t('auto.List.8.f0213e')" min-width="150" show-overflow-tooltip />
            <el-table-column ::label="$t('auto.List.9.d1ef7e')" min-width="180" show-overflow-tooltip>
              <template #default="{ row }">{{ row.description ?? '-' }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.10.ebee05')" width="110">
              <template #default="{ row }">
                <el-tag size="small" effect="plain">{{ channelLabel(row.channel) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.11.3fea7c')" width="100">
              <template #default="{ row }">
                <el-tag :type="pipeStatusType(row.status)" size="small">{{ pipeStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.12.fe2df0')" prop="version" width="80" align="center" />
            <el-table-column :label="$t('auto.List.13.a94e44')" width="160" align="center">
              <template #default="{ row }">
                <span>{{ row.total_runs || 0 }} / {{ row.total_success || 0 }} / {{ row.total_failure || 0 }}</span>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.14.eca37c')" width="170">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.15.2b6bc0')" width="340" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="viewPipe(row)">{{ $t('auto.List.78.607e7a') }}</el-button>
                <el-button v-if="row.status === 'active'" link type="warning" @click="pausePipe(row)">{{ $t('auto.List.79.8d63ef') }}</el-button>
                <el-button v-if="row.status === 'paused'" link type="success" @click="resumePipe(row)">{{ $t('auto.List.80.c7db6d') }}</el-button>
                <el-button v-if="row.status !== 'archived'" link type="info" @click="archivePipe(row)">{{ $t('auto.List.81.2f51c1') }}</el-button>
                <el-button link type="danger" @click="deletePipe(row)">{{ $t('auto.List.82.2f4aad') }}</el-button>
              </template>
            </el-table-column>
            <template #empty><el-empty :description="$t('auto.List.16.b778c8')" /></template>
          </el-table>

          <div class="pager">
            <el-pagination
              background
              layout="total, sizes, prev, pager, next, jumper"
              :total="pipePager.total"
              :page-size="pipePager.page_size"
              :current-page="pipePager.page"
              :page-sizes="[10, 20, 50]"
              @size-change="(s) => { pipePager.page_size = s; loadPipelines() }"
              @current-change="(p) => { pipePager.page = p; loadPipelines() }"
            />
          </div>
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.List.17.b3bfcc')" name="jobs">
        <el-card>
          <div class="toolbar">
            <div class="toolbar-left">
              <el-select
                v-model="jobFilter.pipeline_id"
                :placeholder="$t('auto.List.18.e11f78')"
                clearable
                filterable
                style="width: 220px"
              >
                <el-option
                  v-for="p in pipelines"
                  :key="p.id"
                  :label="`${p.name} (${channelLabel(p.channel)})`"
                  :value="p.id"
                />
              </el-select>
              <el-select
                v-model="jobFilter.channel"
                :placeholder="$t('auto.List.10.ebee05')"
                clearable
                style="width: 140px; margin-left: 8px"
                @change="onJobFilterChange"
              >
                <el-option v-for="c in CHANNEL_OPTIONS" :key="c.value" :label="c.label" :value="c.value" />
              </el-select>
              <el-select
                v-model="jobFilter.state"
                :placeholder="$t('auto.List.11.3fea7c')"
                clearable
                style="width: 140px; margin-left: 8px"
                @change="onJobFilterChange"
              >
                <el-option v-for="s in jobStateOptions" :key="s.value" :label="s.label" :value="s.value" />
              </el-select>
              <el-button style="margin-left: 8px" @click="resetJobFilter">{{ $t('auto.List.83.4b9c32') }}</el-button>
            </div>
            <el-button type="primary" @click="openEnqueueJob">
              <el-icon><Plus /></el-icon>&nbsp;加入任务
            </el-button>
          </div>

          <el-table :data="filteredJobs" v-loading="jobLoading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column label="Pipeline" width="160" show-overflow-tooltip>
              <template #default="{ row }">
                <span>{{ pipelineName(row.pipeline_id) }}</span>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.19.ce8e69')" min-width="140" show-overflow-tooltip>
              <template #default="{ row }">{{ row.customer_id }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.10.ebee05')" width="110">
              <template #default="{ row }">
                <el-tag size="small" effect="plain">{{ channelLabel(row.channel) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.11.3fea7c')" width="110">
              <template #default="{ row }">
                <el-tag :type="jobStateType(row.state)" size="small">{{ jobStateText(row.state) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.20.d5f413')" width="100" align="center">
              <template #default="{ row }">{{ row.retry_count || 0 }}/{{ row.max_retry || 0 }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.21.39f137')" width="100" align="center">
              <template #default="{ row }">{{ row.duration_ms ? row.duration_ms + 'ms' : '-' }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.14.eca37c')" width="170">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.15.2b6bc0')" width="260" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="viewJob(row)">{{ $t('auto.List.78.607e7a') }}</el-button>
                <el-button v-if="canExecute(row.state)" link type="success" @click="executeJob(row)">{{ $t('auto.List.84.1a6aa2') }}</el-button>
                <el-button v-if="canCancel(row.state)" link type="warning" @click="cancelJob(row)">{{ $t('auto.List.85.625fb2') }}</el-button>
                <el-button v-if="row.state === 'failed'" link type="primary" @click="retryJob(row)">{{ $t('auto.List.66.132c5c') }}</el-button>
              </template>
            </el-table-column>
            <template #empty><el-empty :description="$t('auto.List.22.358713')" /></template>
          </el-table>

          <div class="pager">
            <el-pagination
              background
              layout="total, sizes, prev, pager, next, jumper"
              :total="jobPager.total"
              :page-size="jobPager.page_size"
              :current-page="jobPager.page"
              :page-sizes="[10, 20, 50]"
              @size-change="(s) => { jobPager.page_size = s; loadJobs() }"
              @current-change="(p) => { jobPager.page = p; loadJobs() }"
            />
          </div>
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.List.23.5c70d8')" name="stats">
        <el-row :gutter="16" class="stats-row">
          <el-col :xs="12" :sm="8" :md="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-icon stat-icon-green"><el-icon><CircleCheck /></el-icon></div>
              <div class="stat-body">
                <div class="stat-title">{{ $t('auto.List.29.b66435') }}</div>
                <div class="stat-value">{{ overallSuccessRate }}%</div>
                <div class="stat-sub">成功 {{ stats.success || 0 }} / 总计 {{ stats.jobs || 0 }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :xs="12" :sm="8" :md="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-icon stat-icon-red"><el-icon><CircleClose /></el-icon></div>
              <div class="stat-body">
                <div class="stat-title">{{ $t('auto.List.86.b3fe67') }}</div>
                <div class="stat-value">{{ overallFailureRate }}%</div>
                <div class="stat-sub">失败 {{ stats.failed || 0 }} / 总计 {{ stats.jobs || 0 }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :xs="12" :sm="8" :md="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-icon stat-icon-purple"><el-icon><Timer /></el-icon></div>
              <div class="stat-body">
                <div class="stat-title">{{ $t('auto.List.87.b2f29c') }}</div>
                <div class="stat-value">{{ avgDuration }}</div>
                <div class="stat-sub">{{ $t('auto.List.88.57ebd2') }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :xs="12" :sm="8" :md="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-icon stat-icon-orange"><el-icon><Warning /></el-icon></div>
              <div class="stat-body">
                <div class="stat-title">{{ $t('auto.List.89.7a955b') }}</div>
                <div class="stat-value">{{ stats.rate_limited || 0 }}</div>
                <div class="stat-sub">已取消 {{ stats.canceled || 0 }}</div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.90.f1ab24') }}</span>
              <span class="muted">基于近期任务样本（{{ statsSample.length }} 条）</span>
            </div>
          </template>
          <el-table :data="channelStats" v-loading="statsLoading" stripe>
            <el-table-column :label="$t('auto.List.10.ebee05')" width="140">
              <template #default="{ row }">
                <el-tag size="small" effect="plain">{{ channelLabel(row.channel) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="total" :label="$t('auto.List.24.b0677b')" width="120" align="center" />
            <el-table-column prop="success" :label="$t('auto.List.25.330363')" width="100" align="center">
              <template #default="{ row }">
                <span style="color: #10B981">{{ row.success }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="failed" :label="$t('auto.List.26.acd5cb')" width="100" align="center">
              <template #default="{ row }">
                <span style="color: #EF4444">{{ row.failed }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="running" :label="$t('auto.List.27.d679ae')" width="100" align="center" />
            <el-table-column prop="pending" :label="$t('auto.List.28.6139a6')" width="100" align="center" />
            <el-table-column :label="$t('auto.List.29.b66435')" width="120" align="center">
              <template #default="{ row }">{{ row.successRate }}%</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.30.0d652d')" width="140" align="center">
              <template #default="{ row }">{{ row.avgDuration }}</template>
            </el-table-column>
            <template #empty><el-empty :description="$t('auto.List.31.f541ef')" /></template>
          </el-table>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    
    <el-dialog
      v-model="pipeDialogVisible"
      :title="pipeDialogTitle"
      width="780px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <el-form :model="pipeForm" :rules="pipeRules" ref="pipeFormRef" label-width="120px">
        <el-form-item :label="$t('auto.List.32.d7ec2d')" prop="name">
          <el-input v-model="pipeForm.name" :placeholder="$t('auto.List.33.cc2a84')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.34.3bdd08')">
          <el-input v-model="pipeForm.description" type="textarea" :rows="2" :placeholder="$t('auto.List.35.c20cba')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.10.ebee05')" prop="channel">
          <el-select v-model="pipeForm.channel" style="width: 100%" :placeholder="$t('auto.List.36.f46a3f')">
            <el-option v-for="c in CHANNEL_OPTIONS" :key="c.value" :label="c.label" :value="c.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.37.fcebca')" prop="steps">
          <el-select v-model="pipeForm.steps" multiple style="width: 100%" :placeholder="$t('auto.List.38.8fc355')">
            <el-option v-for="s in stepOptions" :key="s.value" :label="s.label" :value="s.value" />
          </el-select>
        </el-form-item>

        <el-divider content-position="left">{{ $t('auto.List.54.90adcc') }}</el-divider>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item :label="$t('auto.List.39.199e1d')">
              <el-input-number v-model="pipeForm.retry_policy.max_retries" :min="0" :max="20" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.List.40.62cfde')">
              <el-input-number v-model="pipeForm.retry_policy.interval_ms" :min="0" :step="500" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.List.41.cf646b')">
              <el-select v-model="pipeForm.retry_policy.backoff" style="width: 100%">
                <el-option :label="$t('auto.List.42.ebdaa1')" value="fixed" />
                <el-option :label="$t('auto.List.43.5f6382')" value="exponential" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.List.44.ee001e')">
              <el-input-number v-model="pipeForm.retry_policy.max_interval_ms" :min="0" :step="1000" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">{{ $t('auto.List.55.f98813') }}</el-divider>
        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item label="QPS">
              <el-input-number v-model="pipeForm.rate_limit.qps" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.List.45.e820ed')">
              <el-input-number v-model="pipeForm.rate_limit.burst" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.List.46.34d1f1')">
              <el-input-number v-model="pipeForm.rate_limit.daily_quota" :min="0" :step="100" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.List.47.35767b')">
              <el-input-number v-model="pipeForm.rate_limit.per_user_limit" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.List.48.178fa2')">
              <el-input-number v-model="pipeForm.rate_limit.cooldown_secs" :min="0" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="pipeDialogVisible = false">{{ $t('auto.List.85.625fb2') }}</el-button>
        <el-button type="primary" :loading="pipeSubmitting" @click="submitPipe">{{ $t('auto.List.91.38cf16') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="pipeDetailVisible" :title="$t('auto.List.49.dacd53')" width="820px">
      <el-descriptions v-if="currentPipe" :column="2" border>
        <el-descriptions-item label="ID">{{ currentPipe.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.32.d7ec2d')">{{ currentPipe.name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.10.ebee05')">{{ channelLabel(currentPipe.channel) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.11.3fea7c')">
          <el-tag :type="pipeStatusType(currentPipe.status)" size="small">{{ pipeStatusText(currentPipe.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.12.fe2df0')">{{ currentPipe.version }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.14.eca37c')">{{ formatTime(currentPipe.created_at) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.34.3bdd08')" :span="2">{{ currentPipe.description || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.50.7ab7e7')">{{ currentPipe.total_runs || 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.51.e7d83a')">{{ currentPipe.total_success || 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.52.d3e480')">{{ currentPipe.total_failure || 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.53.1724cf')">{{ pipeCompletionRate(currentPipe) }}%</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.37.fcebca')" :span="2">
          <div class="steps-list">
            <el-tag
              v-for="(s, i) in pipeStepsList(currentPipe)"
              :key="i"
              class="step-tag"
              :type="s === 'send' ? 'success' : 'info'"
              effect="plain"
            >
              {{ i + 1 }}. {{ stepLabel(s) }}
            </el-tag>
            <span v-if="!pipeStepsList(currentPipe).length" class="muted">-</span>
          </div>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.54.90adcc')" :span="2">
          <pre class="json-pre">{{ formatJSON(currentPipe.retry_policy) }}</pre>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.55.f98813')" :span="2">
          <pre class="json-pre">{{ formatJSON(currentPipe.rate_limit) }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    
    <el-dialog
      v-model="jobDialogVisible"
      :title="$t('auto.List.56.c2c295')"
      width="640px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <el-form :model="jobForm" :rules="jobRules" ref="jobFormRef" label-width="120px">
        <el-form-item label="Pipeline" prop="pipeline_id">
          <el-select
            v-model="jobForm.pipeline_id"
            :placeholder="$t('auto.List.57.e541d3')"
            filterable
            style="width: 100%"
            @change="onJobPipelineChange"
          >
            <el-option
              v-for="p in activePipelines"
              :key="p.id"
              :label="`${p.name} (${channelLabel(p.channel)})`"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.10.ebee05')">
          <el-select v-model="jobForm.channel" :placeholder="$t('auto.List.58.f28b4b')" clearable style="width: 100%">
            <el-option v-for="c in CHANNEL_OPTIONS" :key="c.value" :label="c.label" :value="c.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.59.4af3bd')" prop="customer_id">
          <el-input v-model="jobForm.customer_id" :placeholder="$t('auto.List.60.ae1405')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.61.14d171')">
          <el-input v-model="jobForm.account_id" :placeholder="$t('auto.List.62.da9f21')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.63.f1a93e')">
          <el-input-number v-model="jobForm.max_retry" :min="0" :max="20" />
          <span class="muted" style="margin-left: 8px">{{ $t('auto.List.92.9cc048') }}</span>
        </el-form-item>
        <el-form-item label="Payload" prop="payloadStr">
          <el-input v-model="jobForm.payloadStr" type="textarea" :rows="4" placeholder='例如: {"template":"welcome"}' />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="jobDialogVisible = false">{{ $t('auto.List.85.625fb2') }}</el-button>
        <el-button type="primary" :loading="jobSubmitting" @click="submitJob">{{ $t('auto.List.93.65d1e7') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="jobDetailVisible" :title="$t('auto.List.64.4a98bf')" width="860px">
      <el-descriptions v-if="currentJob" :column="2" border>
        <el-descriptions-item label="ID">{{ currentJob.id }}</el-descriptions-item>
        <el-descriptions-item label="Pipeline">{{ pipelineName(currentJob.pipeline_id) }} (#{{ currentJob.pipeline_id }})</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.10.ebee05')">{{ channelLabel(currentJob.channel) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.11.3fea7c')">
          <el-tag :type="jobStateType(currentJob.state)" size="small">{{ jobStateText(currentJob.state) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.19.ce8e69')">{{ currentJob.customer_id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.65.7035c6')">{{ currentJob.account_id || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.66.132c5c')">{{ currentJob.retry_count || 0 }}/{{ currentJob.max_retry || 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.21.39f137')">{{ currentJob.duration_ms ? currentJob.duration_ms + 'ms' : '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.67.592c59')">{{ formatTime(currentJob.started_at) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.68.7cbadb')">{{ formatTime(currentJob.completed_at) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.69.6cac64')">{{ formatTime(currentJob.next_run_at) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.14.eca37c')">{{ formatTime(currentJob.created_at) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.70.4604d5')" :span="2">
          <span :style="{ color: currentJob.error_message ? '#EF4444' : '' }">{{ currentJob.error_message || '-' }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="Payload" :span="2">
          <pre class="json-pre">{{ formatJSON(currentJob.payload) }}</pre>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.71.8b400a')" :span="2">
          <el-table :data="jobStepResults(currentJob)" size="small" border>
            <el-table-column label="#" type="index" width="50" align="center" />
            <el-table-column :label="$t('auto.List.72.52b365')" width="150">
              <template #default="{ row }">{{ stepLabel(row.step) }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.73.5ad7f5')" width="90" align="center">
              <template #default="{ row }">
                <el-tag :type="row.success ? 'success' : 'danger'" size="small">{{ row.success ? '成功' : '失败' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.21.39f137')" width="90" align="center">
              <template #default="{ row }">{{ row.duration_ms || 0 }}ms</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.74.ad3f65')" min-width="220">
              <template #default="{ row }">
                <span v-if="row.error" style="color: #EF4444">{{ row.error }}</span>
                <pre v-else class="json-pre small">{{ formatJSON(row.output) }}</pre>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!jobStepResults(currentJob).length" :description="$t('auto.List.75.92328f')" :image-size="60" />
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Connection, VideoPlay, VideoPause, FolderOpened, List,
  CircleCheck, CircleClose, Timer, Warning, Share
} from '@element-plus/icons-vue'
import { reachPipelineApi } from '@/api/reachPipeline.js'
import { CHANNEL_OPTIONS, getChannelLabel } from '@/constants/channel'

const stepOptions = [
  {
    value: 'audience',
    label: t('auto.List.94.2c1c8b')
  },
  { value: 'content_prepare', label: t('auto.List.95.dfabcd') },
  { value: 'account_select', label: t('auto.List.96.96b558') },
  { value: 'rate_limit', label: t('auto.List.97.64b702') },
  { value: 'message_gen', label: t('auto.List.98.cb6057') },
  { value: 'send', label: t('auto.List.99.c45a76') },
  { value: 'track_result', label: t('auto.List.100.e45661') },
  { value: 'retry', label: t('auto.List.101.37cd80') },
  { value: 'report', label: t('auto.List.102.8468a0') }
];

const jobStateOptions = [
  { value: 'pending', label: t('auto.List.103.0b5880') },
  { value: 'running', label: t('auto.List.104.46e386') },
  { value: 'success', label: t('auto.List.25.330363') },
  { value: 'failed', label: t('auto.List.26.acd5cb') },
  { value: 'canceled', label: t('auto.List.85.625fb2') },
  { value: 'retrying', label: t('auto.List.105.d7a2a2') },
  { value: 'rate_limited', label: t('auto.List.106.a8ac20') }
]

const channelLabel = getChannelLabel
const stepLabel = (v) => (stepOptions.find((s) => s.value === v) || {}).label || v
const pipeStatusType = (s) => ({ active: 'success', paused: 'warning', archived: 'info' }[s] || 'info')
const pipeStatusText = (s) => ({ active: t('auto.List.107.fe32de'), paused: t('auto.List.79.8d63ef'), archived: t('auto.List.81.2f51c1') }[s] || s)
const jobStateType = (s) => ({
  pending: 'info',
  running: 'warning',
  success: 'success',
  failed: 'danger',
  canceled: 'info',
  retrying: 'warning',
  rate_limited: 'danger'
}[s] || 'info')
const jobStateText = (s) => ({
  pending: t('auto.List.103.0b5880'),
  running: t('auto.List.104.46e386'),
  success: t('auto.List.25.330363'),
  failed: t('auto.List.26.acd5cb'),
  canceled: t('auto.List.85.625fb2'),
  retrying: t('auto.List.105.d7a2a2'),
  rate_limited: t('auto.List.106.a8ac20')
}[s] || s)

const formatTime = (t) => {
  if (!t) return '-'
  const d = new Date(t)
  if (isNaN(d.getTime())) return '-'
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
}
const formatJSON = (obj) => {
  if (obj === null || obj === undefined || obj === '') return '-'
  try {
    return JSON.stringify(obj, null, 2)
  } catch (e) {
    return String(obj)
  }
}

const { t } = useI18n()

const stats = ref({});
const loadStats = async () => {
  try {
    const res = await reachPipelineApi.getStats()
    stats.value = res || {}
  } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}
const archivedCount = computed(() => {
  const total = Number(stats.value.total || 0)
  const active = Number(stats.value.active || 0)
  const paused = Number(stats.value.paused || 0)
  const archived = total - active - paused
  return archived > 0 ? archived : 0
})
const overallSuccessRate = computed(() => {
  const total = Number(stats.value.jobs || 0)
  const succ = Number(stats.value.success || 0)
  if (!total) return 0
  return ((succ / total) * 100).toFixed(1)
})
const overallFailureRate = computed(() => {
  const total = Number(stats.value.jobs || 0)
  const fail = Number(stats.value.failed || 0)
  if (!total) return 0
  return ((fail / total) * 100).toFixed(1)
})

const activeTab = ref('pipelines');
const handleTabChange = (name) => {
  if (name === 'jobs' && !jobs.value.length) loadJobs()
  if (name === 'stats') loadStatsSample()
}

const pipelines = ref([]);
const pipeLoading = ref(false)
const pipeFilter = reactive({ keyword: '', channel: '', status: '' })
const pipePager = reactive({ page: 1, page_size: 20, total: 0 })

const loadPipelines = async () => {
  pipeLoading.value = true
  try {
    const res = await reachPipelineApi.getPipelines({
      page: pipePager.page,
      page_size: pipePager.page_size,
      channel: pipeFilter.channel,
      status: pipeFilter.status
    })
    pipelines.value = res.list || []
    pipePager.total = res.total || 0
  } catch (e) {
    ElMessage.error(e.message || t('auto.List.108.075a63'))
  } finally {
    pipeLoading.value = false
  }
}
const filteredPipelines = computed(() => {
  const kw = (pipeFilter.keyword || '').trim().toLowerCase()
  if (!kw) return pipelines.value
  return pipelines.value.filter((p) => (p.name || '').toLowerCase().includes(kw))
});
const onPipeFilterChange = () => {
  pipePager.page = 1
  loadPipelines()
}
const resetPipeFilter = () => {
  pipeFilter.keyword = ''
  pipeFilter.channel = ''
  pipeFilter.status = ''
  pipePager.page = 1
  loadPipelines()
}
const activePipelines = computed(() => pipelines.value.filter((p) => p.status === 'active'))
const pipelineName = (id) => {
  const p = pipelines.value.find((x) => x.id === id)
  return p ? p.name : `#${id}`
}

const pipeDialogVisible = ref(false);
const pipeDialogTitle = ref(t('auto.List.109.3d5a23'))
const pipeFormRef = ref()
const pipeSubmitting = ref(false)

const defaultPipeForm = () => ({
  name: '',
  description: '',
  channel: 'wecom',
  steps: stepOptions.map((s) => s.value),
  retry_policy: { max_retries: 3, interval_ms: 1000, backoff: 'exponential', max_interval_ms: 60000 },
  rate_limit: { qps: 10, burst: 20, daily_quota: 10000, per_user_limit: 3, cooldown_secs: 60 }
})
const pipeForm = reactive(defaultPipeForm())
const pipeRules = {
  name: [{ required: true, message: i18n.global.t(t('auto.List.110.06e2f8')), trigger: 'blur' }],
  channel: [{ required: true, message: i18n.global.t(t('auto.List.36.f46a3f')), trigger: 'change' }],
  steps: [{ required: true, type: 'array', min: 1, message: i18n.global.t(t('auto.List.111.536ab2')), trigger: 'change' }]
}

const openCreatePipe = () => {
  Object.assign(pipeForm, defaultPipeForm())
  pipeDialogTitle.value = t('auto.List.109.3d5a23')
  pipeDialogVisible.value = true
}
const submitPipe = async () => {
  if (!pipeFormRef.value) return
  await pipeFormRef.value.validate(async (valid) => {
    if (!valid) return
    if (!pipeForm.steps.includes('send')) {
      ElMessage.warning(i18n.global.t(t('auto.List.112.62c685')))
      return
    }
    pipeSubmitting.value = true
    try {
      const data = {
        name: pipeForm.name,
        description: pipeForm.description,
        channel: pipeForm.channel,
        steps: pipeForm.steps,
        retry_policy: pipeForm.retry_policy,
        rate_limit: pipeForm.rate_limit
      }
      await reachPipelineApi.createPipeline(data)
      ElMessage.success(i18n.global.t(t('auto.List.113.04a691')))
      pipeDialogVisible.value = false
      loadPipelines()
      loadStats()
    } catch (e) {
      ElMessage.error(e.message || t('auto.List.114.a88928'))
    } finally {
      pipeSubmitting.value = false
    }
  })
}

const pausePipe = async (row) => {
  try {
    await reachPipelineApi.pausePipeline(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.115.a2d930')))
    loadPipelines()
    loadStats()
  } catch (e) {
    ElMessage.error(e.message || t('auto.List.116.5fa802'))
  }
}
const resumePipe = async (row) => {
  try {
    await reachPipelineApi.resumePipeline(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.117.b70e8e')))
    loadPipelines()
    loadStats()
  } catch (e) {
    ElMessage.error(e.message || t('auto.List.116.5fa802'))
  }
}
const archivePipe = async (row) => {
  try {
    await ElMessageBox.confirm(`确定归档 "${row.name}"？归档后不可再加入任务`, t('auto.List.118.e83a25'), { type: 'warning' })
    await reachPipelineApi.archivePipeline(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.119.c3ba16')))
    loadPipelines()
    loadStats()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message || t('auto.List.116.5fa802'))
  }
}
const deletePipe = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除 "${row.name}"？此操作不可恢复`, t('auto.List.118.e83a25'), { type: 'warning' })
    await reachPipelineApi.deletePipeline(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.120.5cc232')))
    loadPipelines()
    loadStats()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message || t('auto.List.116.5fa802'))
  }
}

const resetRateLimit = async () => {
  if (!pipeFilter.channel) {
    ElMessage.warning(i18n.global.t(t('auto.List.121.df3b6d')))
    return
  }
  try {
    await ElMessageBox.confirm(`确定重置渠道 "${channelLabel(pipeFilter.channel)}" 的限流状态？`, t('auto.List.118.e83a25'), { type: 'warning' })
    await reachPipelineApi.resetRateLimit(pipeFilter.channel)
    ElMessage.success(i18n.global.t(t('auto.List.122.6fb95c')))
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message || t('auto.List.116.5fa802'))
  }
};

const pipeDetailVisible = ref(false);
const currentPipe = ref(null)
const viewPipe = async (row) => {
  try {
    const res = await reachPipelineApi.getPipeline(row.id)
    currentPipe.value = res
    pipeDetailVisible.value = true
  } catch (e) {
    ElMessage.error(e.message || t('auto.List.123.181421'))
  }
}
const pipeStepsList = (p) => (Array.isArray(p && p.steps) ? p.steps.map((s) => String(s)) : [])
const pipeCompletionRate = (p) => {
  const runs = Number(p && p.total_runs || 0)
  const succ = Number(p && p.total_success || 0)
  if (!runs) return 0
  return ((succ / runs) * 100).toFixed(1)
}

const jobs = ref([]);
const jobLoading = ref(false)
const jobFilter = reactive({ pipeline_id: null, channel: '', state: '' })
const jobPager = reactive({ page: 1, page_size: 20, total: 0 })

const loadJobs = async () => {
  jobLoading.value = true
  try {
    const res = await reachPipelineApi.getJobs({
      page: jobPager.page,
      page_size: jobPager.page_size,
      channel: jobFilter.channel,
      state: jobFilter.state
    })
    jobs.value = res.list || []
    jobPager.total = res.total || 0
  } catch (e) {
    ElMessage.error(e.message || t('auto.List.124.f71ace'))
  } finally {
    jobLoading.value = false
  }
}
const filteredJobs = computed(() => {
  if (!jobFilter.pipeline_id) return jobs.value
  return jobs.value.filter((j) => j.pipeline_id === jobFilter.pipeline_id)
});
const onJobFilterChange = () => {
  jobPager.page = 1
  loadJobs()
}
const resetJobFilter = () => {
  jobFilter.pipeline_id = null
  jobFilter.channel = ''
  jobFilter.state = ''
  jobPager.page = 1
  loadJobs()
}

const canExecute = (s) => ['pending', 'retrying', 'rate_limited'].includes(s)
const canCancel = (s) => ['pending', 'running', 'retrying', 'rate_limited'].includes(s)

const executeJob = async (row) => {
  try {
    await reachPipelineApi.executeJob(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.125.c044a1')))
    loadJobs()
    loadStats()
  } catch (e) {
    ElMessage.error(e.message || t('auto.List.126.1c83d7'))
    loadJobs()
  }
}
const cancelJob = async (row) => {
  try {
    await ElMessageBox.confirm(`确定取消任务 #${row.id}？`, t('auto.List.118.e83a25'), { type: 'warning' })
    await reachPipelineApi.cancelJob(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.127.2111cc')))
    loadJobs()
    loadStats()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message || t('auto.List.116.5fa802'))
  }
}
const retryJob = async (row) => {
  try {
    await reachPipelineApi.retryJob(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.128.60c37f')))
    loadJobs()
    loadStats()
  } catch (e) {
    ElMessage.error(e.message || t('auto.List.116.5fa802'))
  }
}

const jobDialogVisible = ref(false);
const jobFormRef = ref()
const jobSubmitting = ref(false)
const jobForm = reactive({ pipeline_id: null, channel: '', customer_id: '', account_id: '', max_retry: 0, payloadStr: '{}' })
const jobRules = {
  pipeline_id: [{ required: true, message: i18n.global.t(t('auto.List.129.d289d0')), trigger: 'change' }],
  customer_id: [{ required: true, message: i18n.global.t(t('auto.List.130.fe50d1')), trigger: 'blur' }],
  payloadStr: [{ required: true, message: i18n.global.t(t('auto.List.131.968397')), trigger: 'blur' }]
}
const openEnqueueJob = () => {
  Object.assign(jobForm, { pipeline_id: null, channel: '', customer_id: '', account_id: '', max_retry: 0, payloadStr: '{}' })
  jobDialogVisible.value = true
}
const onJobPipelineChange = (id) => {
  const p = pipelines.value.find((x) => x.id === id)
  if (p && !jobForm.channel) jobForm.channel = p.channel
}
const submitJob = async () => {
  if (!jobFormRef.value) return
  await jobFormRef.value.validate(async (valid) => {
    if (!valid) return
    let payload
    try {
      payload = JSON.parse(jobForm.payloadStr)
    } catch (e) {
      ElMessage.error(i18n.global.t('Payload 不是有效的 JSON'))
      return
    }
    jobSubmitting.value = true
    try {
      const data = {
        pipeline_id: jobForm.pipeline_id,
        channel: jobForm.channel,
        customer_id: jobForm.customer_id,
        account_id: jobForm.account_id,
        max_retry: jobForm.max_retry,
        payload
      }
      await reachPipelineApi.enqueueJob(data)
      ElMessage.success(i18n.global.t(t('auto.List.132.8ad081')))
      jobDialogVisible.value = false
      activeTab.value = 'jobs'
      jobPager.page = 1
      loadJobs()
      loadStats()
    } catch (e) {
      ElMessage.error(e.message || t('auto.List.133.ca16d3'))
    } finally {
      jobSubmitting.value = false
    }
  })
}

const jobDetailVisible = ref(false);
const currentJob = ref(null)
const viewJob = async (row) => {
  try {
    const res = await reachPipelineApi.getJob(row.id)
    currentJob.value = res
    jobDetailVisible.value = true
  } catch (e) {
    ElMessage.error(e.message || t('auto.List.123.181421'))
  }
}
const jobStepResults = (j) => (Array.isArray(j && j.step_results) ? j.step_results : [])

const statsSample = ref([]);
const statsLoading = ref(false)
const loadStatsSample = async () => {
  statsLoading.value = true
  try {
    const res = await reachPipelineApi.getJobs({ page: 1, page_size: 200 });
    statsSample.value = res.list || []
  } catch (e) {
    ElMessage.error(e.message || t('auto.List.134.f8b8e1'))
  } finally {
    statsLoading.value = false
  }
}
const avgDuration = computed(() => {
  const valid = statsSample.value.filter((j) => Number(j.duration_ms) > 0)
  if (!valid.length) return '-'
  const avg = valid.reduce((s, j) => s + Number(j.duration_ms), 0) / valid.length
  return avg < 1000 ? `${Math.round(avg)}ms` : `${(avg / 1000).toFixed(2)}s`
})
const channelStats = computed(() => {
  const map = {}
  statsSample.value.forEach((j) => {
    const ch = j.channel || 'unknown'
    if (!map[ch]) {
      map[ch] = { channel: ch, total: 0, success: 0, failed: 0, running: 0, pending: 0, durationSum: 0, durationCnt: 0 }
    }
    const item = map[ch]
    item.total++
    if (j.state === 'success') item.success++
    else if (j.state === 'failed') item.failed++
    else if (j.state === 'running') item.running++
    else if (j.state === 'pending') item.pending++
    if (Number(j.duration_ms) > 0) {
      item.durationSum += Number(j.duration_ms)
      item.durationCnt++
    }
  })
  return Object.values(map).map((item) => {
    const rate = item.total ? ((item.success / item.total) * 100).toFixed(1) : '0.0'
    const avg = item.durationCnt ? Math.round(item.durationSum / item.durationCnt) : 0
    return {
      ...item,
      successRate: rate,
      avgDuration: avg < 1000 ? `${avg}ms` : `${(avg / 1000).toFixed(2)}s`
    }
  })
})

onMounted(() => {
  loadStats()
  loadPipelines()
  loadJobs()
});
</script>

<style scoped lang="scss">
.reach-pipeline-page {
  padding: 20px;
}

.rp-tabs {
  :deep(.el-tabs__content) {
    padding-top: 8px;
  }
}

.stats-row {
  margin-bottom: 16px;
}

.stat-card {
  margin-bottom: 12px;
  :deep(.el-card__body) {
    display: flex;
    align-items: center;
    padding: 18px 20px;
  }
  .stat-icon {
    width: 52px;
    height: 52px;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 16px;
    font-size: 26px;
    color: #fff;
    flex-shrink: 0;
  }
  .stat-icon-blue { background: linear-gradient(135deg, #4F46E5, #66b1ff); }
  .stat-icon-green { background: linear-gradient(135deg, #10B981, #95d475); }
  .stat-icon-green2 { background: linear-gradient(135deg, #529b2e, #10B981); }
  .stat-icon-orange { background: linear-gradient(135deg, #F59E0B, #f0b87a); }
  .stat-icon-red { background: linear-gradient(135deg, #EF4444, #f89898); }
  .stat-icon-purple { background: linear-gradient(135deg, #9c6cd6, #b589e0); }
  .stat-icon-cyan { background: linear-gradient(135deg, #13c2c2, #5cdbd3); }

  .stat-body {
    flex: 1;
    min-width: 0;
  }
  .stat-title {
    font-size: 13px;
    color: #909399;
    margin-bottom: 6px;
  }
  .stat-value {
    font-size: 26px;
    font-weight: 600;
    color: #303133;
    line-height: 1.2;
  }
  .stat-sub {
    font-size: 12px;
    color: #c0c4cc;
    margin-top: 4px;
  }
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 8px;
}
.toolbar-left {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.steps-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.step-tag {
  margin: 0;
}

.muted {
  color: #c0c4cc;
  font-size: 13px;
}

.json-pre {
  margin: 0;
  padding: 8px 10px;
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 12px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 240px;
  overflow: auto;
}
.json-pre.small {
  padding: 4px 6px;
  max-height: 120px;
}
</style>
