<template>
  <div class="guide">
    <p class="lead">{{ $t('auto.HostInstallGuide.973ffb79') }}</p>
    <ol class="steps">
      <li>{{ $t('auto.HostInstallGuide.03f04f0e') }}<code>cd user-web/browser_automation && npm install && npm run build</code>
        <div class="note">{{ $t('auto.HostInstallGuide.2b83f616') }}<code>dist</code>{{ $t('auto.HostInstallGuide.11ae6e1f') }}</div></li>
      <li>chrome://extensions → 打开「开发者模式」→「加载已解压的扩展程序」→ 选择
        <code>user-web/browser_automation/dist</code></li>
      <li>{{ $t('auto.HostInstallGuide.51105c7e') }}<code>user-server/cmd/nm-host/install.sh</code>
        <div class="note">{{ $t('auto.HostInstallGuide.d009e596') }}<code>key</code>{{ $t('auto.HostInstallGuide.ab4aa970') }}</div></li>
      <li>{{ $t('auto.HostInstallGuide.bb57d8a9') }}<code>~/.hivemtk/nm_host.conf</code>{{ $t('auto.HostInstallGuide.01d7aa49') }}<code>token=</code> 行
        <div class="note">{{ $t('auto.HostInstallGuide.98300539') }}</div></li>
    </ol>
    <p class="lead">{{ $t('auto.HostInstallGuide.cfe7bf82') }}</p>
  </div>
</template>

<style scoped>
.guide { line-height: 1.8; }
.lead { margin: 0 0 6px; }
.steps { margin: 0 0 6px 18px; padding: 0; }
.steps li { margin-bottom: 6px; }
.note { color: #909399; font-size: 12px; line-height: 1.6; }
</style>
