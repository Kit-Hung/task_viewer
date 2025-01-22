<template>
  <div ref="graphContainer" class="dependency-graph"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import * as d3 from 'd3'

const props = defineProps<{
  tasks: Array<{ id: number; name: string }>
  dependencies: Array<{ source: number; target: number }>
  selectedTasks: number[]
}>()

const emit = defineEmits<{
  (e: 'select', taskIds: number[]): void
  (e: 'node-dblclick', taskId: number): void
}>()

const graphContainer = ref<HTMLElement>()

const renderGraph = () => {
  if (!graphContainer.value) return

  console.log('Rendering graph with data:', {
    tasks: props.tasks,
    dependencies: props.dependencies,
    selectedTasks: props.selectedTasks
  })

  // 清除现有的图表
  d3.select(graphContainer.value).selectAll('*').remove()

  const width = graphContainer.value.clientWidth
  const height = graphContainer.value.clientHeight || 500

  const svg = d3.select(graphContainer.value)
    .append('svg')
    .attr('width', '100%')
    .attr('height', '100%')
    .attr('viewBox', `0 0 ${width} ${height}`)
    .attr('preserveAspectRatio', 'xMidYMid meet')

  const g = svg.append('g')

  // 创建力导向图布局
  const simulation = d3.forceSimulation(props.tasks)
    .force('link', d3.forceLink()
      .id((d: any) => d.id)
      .distance(120))
    .force('charge', d3.forceManyBody().strength(-400))
    .force('center', d3.forceCenter(width / 2, height / 2))
    .force('x', d3.forceX(width / 2).strength(0.05))
    .force('y', d3.forceY(height / 2).strength(0.05))

  // 添加缩放功能
  const zoom = d3.zoom()
    .scaleExtent([0.2, 3])
    .on('zoom', (event) => {
      g.attr('transform', event.transform)
    })

  svg.call(zoom as any)
    .call(zoom.transform as any, d3.zoomIdentity
      .translate(width/2, height/2)
      .scale(0.8)
      .translate(-width/2, -height/2))

  // 添加箭头标记
  svg.append('defs').append('marker')
    .attr('id', 'arrowhead')
    .attr('viewBox', '-0 -5 10 10')
    .attr('refX', 20)
    .attr('refY', 0)
    .attr('orient', 'auto')
    .attr('markerWidth', 6)
    .attr('markerHeight', 6)
    .append('path')
    .attr('d', 'M 0,-4 L 8,0 L 0,4')
    .attr('fill', '#409EFF')

  // 绘制连线
  const links = g.append('g')
    .selectAll('line')
    .data(props.dependencies.map(d => ({
      source: props.tasks.find(t => t.id === d.source),
      target: props.tasks.find(t => t.id === d.target)
    })))
    .join('line')
    .attr('stroke', '#409EFF')
    .attr('stroke-width', 1.5)
    .attr('marker-end', 'url(#arrowhead)')

  console.log('Created links:', links.size(), 'from dependencies:', props.dependencies)

  // 创建节点组
  const nodeGroup = g.append('g')
    .selectAll('g')
    .data(props.tasks)
    .enter()
    .append('g')
    .call(d3.drag()
      .on('start', dragstarted)
      .on('drag', dragged)
      .on('end', dragended))

  // 绘制节点圆形背景
  nodeGroup.append('circle')
    .attr('r', 14)
    .attr('fill', (d: any) => props.selectedTasks?.includes(d.id) ? '#409EFF' : '#67C23A')
    .attr('stroke', (d: any) => props.selectedTasks?.includes(d.id) ? '#79bbff' : '#95d475')
    .attr('stroke-width', 1.5)
    .style('cursor', 'pointer')
    .style('filter', 'drop-shadow(0 2px 4px rgba(0,0,0,0.1))')

  // 添加节点文本
  nodeGroup.append('text')
    .text((d: any) => d.name)
    .attr('text-anchor', 'middle')
    .attr('dy', '.35em')
    .attr('fill', '#fff')
    .attr('font-size', '10px')
    .style('pointer-events', 'none')
    .style('text-shadow', '0 1px 2px rgba(0,0,0,0.2)')

  // 更新力导向图
  simulation
    .nodes(props.tasks)
    .on('tick', ticked)

  const linkForce = simulation.force<d3.ForceLink<any, any>>('link')!
  linkForce.links(props.dependencies.map(d => ({
    source: props.tasks.find(t => t.id === d.source),
    target: props.tasks.find(t => t.id === d.target)
  })))

  // 处理节点拖拽
  function dragstarted(event: any) {
    if (!event.active) simulation.alphaTarget(0.3).restart()
    event.subject.fx = event.subject.x
    event.subject.fy = event.subject.y
  }

  function dragged(event: any) {
    event.subject.fx = event.x
    event.subject.fy = event.y
  }

  function dragended(event: any) {
    if (!event.active) simulation.alphaTarget(0)
    event.subject.fx = null
    event.subject.fy = null
  }

  // 更新节点和连线位置
  function ticked() {
    links
      .attr('x1', (d: any) => d.source?.x || 0)
      .attr('y1', (d: any) => d.source?.y || 0)
      .attr('x2', (d: any) => {
        if (!d.target?.x || !d.source?.x) return 0
        const dx = d.target.x - d.source.x
        const dy = d.target.y - d.source.y
        const angle = Math.atan2(dy, dx)
        return d.target.x - 20 * Math.cos(angle)
      })
      .attr('y2', (d: any) => {
        if (!d.target?.y || !d.source?.y) return 0
        const dx = d.target.x - d.source.x
        const dy = d.target.y - d.source.y
        const angle = Math.atan2(dy, dx)
        return d.target.y - 20 * Math.sin(angle)
      })

    nodeGroup.attr('transform', (d: any) => `translate(${d.x},${d.y})`)
  }

  // 添加点击事件
  nodeGroup.on('click', (event: any, d: any) => {
    const taskId = d.id
    const newSelectedTasks = props.selectedTasks?.includes(taskId)
      ? props.selectedTasks.filter(id => id !== taskId)
      : [...(props.selectedTasks || []), taskId]
    emit('select', newSelectedTasks)
  })

  // 添加双击事件
  nodeGroup.on('dblclick', (event: any, d: any) => {
    emit('node-dblclick', d.id)
  })
}

onMounted(() => {
  renderGraph()
})

watch([() => props.tasks, () => props.dependencies, () => props.selectedTasks], () => {
  renderGraph()
})
</script>

<style scoped>
.dependency-graph {
  width: 100%;
  height: 100%;
  min-height: 500px;
}
</style> 