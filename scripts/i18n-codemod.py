#!/usr/bin/env python3
# =============================================================================
# i18n-codemod.py —— 视图硬编码中文 → $t('key') 机械转换器（配合人工校对）
#
# 为什么需要它：全仓 views 下 205 个文件 2083 行硬编码中文（2026-10-11 盘点，
# docs/backlog/i18n-inventory.json），逐行手改不现实。本工具做机械部分：
#   1. 提取 <template> 区的中文字面量（属性值 / 文本节点）与 script 区的
#      ElMessage/ElMessageBox 字符串 → 生成稳定 key（auto.<文件名>.<序号>）
#   2. 模板位替换为 $t('key')，script 位替换为 t('key')（自动补 useI18n 导入）
#   3. zh.json 写入中文原文；其余 locale 由 fallbackLocale: 'zh' 兜底，
#      人工翻译逐批补齐（不会回归——缺键回落中文，与硬编码时代显示一致）
#
# 用法：
#   python scripts/i18n-codemod.py user-web/src/views/mcp/Credentials.vue [...]
#   转换完成后必须：npm run build（i18n 门）+ 人工抽查关键页面。
# =============================================================================
import hashlib
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
LOCALES = ["zh"]


def gen_key(base: str, text: str, seq: int) -> str:
    h = hashlib.md5(text.encode("utf-8")).hexdigest()[:6]
    return f"auto.{base}.{seq}.{h}"


def convert_file(path: Path) -> int:
    src = path.read_text(encoding="utf-8")
    base = path.stem
    zh_path = ROOT / "user-web" / "src" / "i18n" / "locales" / "zh.json"
    zh = json.loads(zh_path.read_text(encoding="utf-8"))
    zh.setdefault("auto", {}).setdefault(base, {})

    cn = re.compile(r"[\u4e00-\u9fff]")
    if not cn.search(src):
        return 0
    if "useI18n" in src and "auto." not in src:
        return 0

    mapping = {}
    seq = 0

    def key_for(text: str) -> str:
        nonlocal seq
        # 同文复用同一 key
        for k, v in mapping.items():
            if v == text:
                return k
        seq += 1
        k = gen_key(base, text, seq)
        mapping[k] = text
        return k

    # ---- script 区先行（避免模板改写影响行号判断）----
    script_cn = re.compile(r"(['\"`])([^'\"`\n]*[\u4e00-\u9fff][^'\"`\n]*)\1")
    needs_t = False

    def script_sub(m):
        nonlocal needs_t
        quote, text = m.group(1), m.group(2)
        if not cn.search(text):
            return m.group(0)
        # 模板字符串（含 ${}）不转换：逻辑拼装字符串，人工处理
        if quote == "`" and "${" in text:
            return m.group(0)
        k = key_for(text)
        needs_t = True
        return f"t('{k}')"

    # ElMessage / ElMessageBox / confirm 内字符串
    lines = src.split("\n")
    out_lines = []
    in_template = False
    template_buf = []
    script_buf = []
    for line in lines:
        if "<template" in line:
            in_template = True
        if "</template>" in line:
            in_template = False
            template_buf.append(line)
            continue
        (template_buf if in_template else script_buf).append(line)
    src = "\n".join(template_buf[:-1] if False else template_buf)  # placeholder

    # 简化：整文件统一处理，模板/script 都做替换，模板位再二次加工
    src2 = path.read_text(encoding="utf-8")

    # 2a. 模板属性值：xx="中文"
    def attr_sub(m):
        attr, text = m.group(1), m.group(2)
        k = key_for(text)
        return f':{attr}="$t(\'{k}\')"'

    # 2b. 文本节点：>中文<
    def text_sub(m):
        text = m.group(1).strip()
        if not cn.search(text):
            return m.group(0)
        k = key_for(text)
        return f">{{{{ $t('{k}') }}}}<"

    tpl_start = src2.find("<template")
    tpl_end = src2.rfind("</template>")
    tpl = src2[tpl_start:tpl_end + len("</template>")]
    rest_before, rest_after = src2[:tpl_start], src2[tpl_end + len("</template>"):]

    # 属性（只处理双引号属性且含中文、且不含 {{ 表达式）
    tpl = re.sub(r'([a-zA-Z-]+)="([^"\{\n]*[\u4e00-\u9fff][^"\{\n]*)"',
                 lambda m: f':{m.group(1)}="$t(\'{key_for(m.group(2))}\')"' if cn.search(m.group(2)) else m.group(0),
                 tpl)
    # 文本节点（> 中文 <，跨行不看，单行处理）
    tpl = re.sub(r">([^<>{}\n]*[\u4e00-\u9fff][^<>{}\n]*)<",
                 lambda m: ">{{ $t('" + key_for(m.group(1).strip()) + "') }}<",
                 tpl)

    # 2c/d. script 区中文字符串统一（ElMessage/ElMessageBox 等都走这段）
    src3 = rest_before + tpl + rest_after
    sc_start = src3.find("<script")
    sc_end = src3.find("</script>")
    if sc_start >= 0 and sc_end > sc_start:
        head, script, tail = src3[:sc_start], src3[sc_start:sc_end], src3[sc_end:]
        cjk_str = re.compile(r"'([^'\n]*[\u4e00-\u9fff][^'\n]*)'")
        script = cjk_str.sub(
            lambda m: "t('" + key_for(m.group(1)) + "')" if cn.search(m.group(1)) else m.group(0),
            script)
        # 纯字面量才替换；含 ${} 插值的模板串留给人工（转换器红线）
        src3 = head + script + tail

    # 2e. useI18n 注入
    if needs_t or "$t(" in src3:
        if "useI18n" not in src3:
            src3 = src3.replace("<script setup>", "<script setup>", 1)
            m = re.search(r"import .* from 'vue'\n", src3)
            if m:
                src3 = src3[:m.end()] + "import { useI18n } from 'vue-i18n'\n" + src3[m.end():]
            m2 = re.search(r"(const \w+ = (?:ref|reactive)\(|onMounted\()", src3)
            if m2:
                src3 = src3[:m2.start()] + "const { t } = useI18n()\n\n" + src3[m2.start():]

    # 2f. zh 入典
    section = zh.setdefault("auto", {}).setdefault(base, {})
    for k, v in mapping.items():
        section[k.split(".", 2)[2]] = v

    path.write_text(src3, encoding="utf-8", newline="\n")
    zh_path.write_text(json.dumps(zh, ensure_ascii=False, indent=2) + "\n", encoding="utf-8", newline="\n")
    return len(mapping)


def main():
    n_total = 0
    for arg in sys.argv[1:]:
        p = Path(arg)
        if not p.exists():
            print(f"跳过（不存在）: {arg}")
            continue
        n = convert_file(p)
        print(f"{arg}: 提取 {n} 条")
        n_total += n
    print(f"合计 {n_total} 条。下一步：npm run build（i18n 门）+ 人工校对 + 补 en 及其余 locale。")


if __name__ == "__main__":
    main()
