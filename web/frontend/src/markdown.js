// AI-assisted code
export function renderMarkdown(src) {
  let html = ''
  const lines = src.split('\n')
  let i = 0

  while (i < lines.length) {
    const line = lines[i]

    if (line.startsWith('```')) {
      i++
      let code = ''
      while (i < lines.length && !lines[i].startsWith('```')) {
        code += esc(lines[i]) + '\n'
        i++
      }
      i++
      html += '<pre>' + code.trimEnd() + '</pre>'
      continue
    }

    const hm = line.match(/^(#{1,4})\s+(.+)/)
    if (hm) {
      html += `<h${hm[1].length}>${inl(hm[2])}</h${hm[1].length}>`
      i++
      continue
    }

    if (line.includes('|') && line.trim().startsWith('|')) {
      const rows = []
      while (i < lines.length && lines[i].includes('|') && lines[i].trim().startsWith('|')) {
        const cells = lines[i].split('|').slice(1, -1).map(c => c.trim())
        if (!cells.every(c => /^[-:]+$/.test(c))) rows.push(cells)
        i++
      }
      if (rows.length) {
        html += '<table>'
        html += '<tr>' + rows[0].map(c => '<th>' + inl(c) + '</th>').join('') + '</tr>'
        for (let r = 1; r < rows.length; r++) {
          html += '<tr>' + rows[r].map(c => '<td>' + inl(c) + '</td>').join('') + '</tr>'
        }
        html += '</table>'
      }
      continue
    }

    if (line.match(/^[-*]\s/)) {
      html += '<ul>'
      while (i < lines.length && lines[i].match(/^[-*]\s/)) {
        html += '<li>' + inl(lines[i].replace(/^[-*]\s/, '')) + '</li>'
        i++
      }
      html += '</ul>'
      continue
    }

    if (!line.trim()) { i++; continue }

    let para = ''
    while (i < lines.length && lines[i].trim() &&
           !lines[i].startsWith('#') && !lines[i].startsWith('```') &&
           !lines[i].match(/^[-*]\s/) &&
           !(lines[i].includes('|') && lines[i].trim().startsWith('|'))) {
      para += (para ? ' ' : '') + lines[i]
      i++
    }
    html += '<p>' + inl(para) + '</p>'
  }

  return html
}

function inl(t) {
  return t
    .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
    .replace(/`([^`]+)`/g, '<code>$1</code>')
    .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
    .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2">$1</a>')
}

function esc(t) {
  return t.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
}
