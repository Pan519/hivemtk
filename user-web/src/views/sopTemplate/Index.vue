<template>
  <div class="sop-template-index-page">
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('auto.Index.b9ce6f82') }}</h2>
          <p class="subtitle">{{ $t('auto.Index.c3f3df50') }}</p>
        </div>
      </div>
    </el-card>

    <el-row :gutter="20">
      <el-col :span="8">
        <el-card shadow="hover" class="entry-card" @click="$router.push('/sop-template/list')">
          <div class="entry-icon"><el-icon><Tickets /></el-icon></div>
          <div class="entry-title">{{ $t('auto.Index.3996b920') }}</div>
          <div class="entry-desc">{{ $t('auto.Index.138ca0f3') }}</div>
          <el-button type="primary" plain>
            进入
            <el-icon class="el-icon--right"><ArrowRight /></el-icon>
          </el-button>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="entry-card">
          <div class="entry-icon"><el-icon><Connection /></el-icon></div>
          <div class="entry-title">{{ $t('auto.Index.c50cce7e') }}</div>
          <div class="entry-desc">{{ $t('auto.Index.126208cb') }}</div>
          <el-button type="primary" plain @click="$router.push('/sop-template/list')">
            打开
            <el-icon class="el-icon--right"><ArrowRight /></el-icon>
          </el-button>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="entry-card">
          <div class="entry-icon"><el-icon><Cpu /></el-icon></div>
          <div class="entry-title">{{ $t('auto.Index.62970284') }}</div>
          <div class="entry-desc">{{ $t('auto.Index.727b3b7c') }}</div>
          <el-button type="primary" plain @click="$router.push('/aiAgent/list')">
            去绑定
            <el-icon class="el-icon--right"><ArrowRight /></el-icon>
          </el-button>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { Tickets, Cpu, Connection, ArrowRight } from '@element-plus/icons-vue'
</script>

<style scoped>
.sop-template-index-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.header-card :deep(.el-card__body) {
  padding: 18px 24px;
}

.header-content h2 {
  margin: 0 0 4px 0;
  font-size: 20px;
  font-weight: 600;
}

.subtitle {
  margin: 0;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.entry-card {
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
  text-align: center;
  padding: 8px 0;
}

.entry-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.entry-icon {
  font-size: 36px;
  color: var(--el-color-primary);
  margin-bottom: 8px;
}

.entry-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 6px;
}

.entry-desc {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.5;
  margin-bottom: 12px;
  min-height: 40px;
}
</style>
