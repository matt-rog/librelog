// AI-assisted code
import architecture from '../../../docs/architecture.md?raw'
import selfHosting from '../../../docs/self-hosting.md?raw'
import configuration from '../../../docs/configuration.md?raw'
import api from '../../../docs/api.md?raw'

export const pages = [
  { slug: 'architecture', title: 'Architecture', content: architecture },
  { slug: 'self-hosting', title: 'Self-Hosting', content: selfHosting },
  { slug: 'configuration', title: 'Configuration', content: configuration },
  { slug: 'api', title: 'API Reference', content: api },
]

export function getPage(slug) {
  return pages.find(p => p.slug === slug)
}
