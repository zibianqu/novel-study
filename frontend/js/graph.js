layui.use(['layer'], function() {
    const layer = layui.layer;

    // 检查登录
    if (!API.getToken()) {
        location.href = 'index.html';
        return;
    }

    const userInfo = API.getUserInfo();
    if (userInfo) {
        document.getElementById('username').textContent = userInfo.username;
    }

    // 退出登录
    document.getElementById('logout').addEventListener('click', function() {
        layer.confirm('确定要退出登录吗？', { icon: 3 }, function() {
            localStorage.removeItem(STORAGE_KEYS.TOKEN);
            localStorage.removeItem(STORAGE_KEYS.USER_INFO);
            location.href = 'index.html';
        });
    });

    // 获取项目 ID
    const projectId = getQueryParam('project') || 1;

    // D3.js 图谱变量
    let svg, simulation, link, node, label;
    const width = document.getElementById('graphContainer').clientWidth;
    const height = document.getElementById('graphContainer').clientHeight;

    // 初始化图谱
    initGraph();
    loadGraph();

    // 刷新按钮
    document.getElementById('refreshBtn').addEventListener('click', loadGraph);

    function initGraph() {
        svg = d3.select('#graphContainer')
            .append('svg')
            .attr('width', width)
            .attr('height', height);

        // 创建力导向图
        simulation = d3.forceSimulation()
            .force('link', d3.forceLink().id(d => d.id).distance(150))
            .force('charge', d3.forceManyBody().strength(-400))
            .force('center', d3.forceCenter(width / 2, height / 2));

        // 添加图例
        addLegend();
    }

    async function loadGraph() {
        try {
            const data = await API.get(`/graph/project/${projectId}`);
            renderGraph(data.nodes || [], data.relations || []);
        } catch (error) {
            layer.msg('加载图谱失败', { icon: 2 });
            console.error(error);
        }
    }

    function renderGraph(nodes, links) {
        // 更新统计
        document.getElementById('nodeCount').textContent = nodes.length;
        document.getElementById('relationCount').textContent = links.length;

        // 清除旧元素
        svg.selectAll('.link').remove();
        svg.selectAll('.node').remove();
        svg.selectAll('.node-label').remove();

        // 绘制连接线
        link = svg.selectAll('.link')
            .data(links)
            .enter().append('line')
            .attr('class', 'link');

        // 绘制节点
        node = svg.selectAll('.node')
            .data(nodes)
            .enter().append('circle')
            .attr('class', d => `node node-${d.type}`)
            .attr('r', 20)
            .call(d3.drag()
                .on('start', dragstarted)
                .on('drag', dragged)
                .on('end', dragended));

        // 添加标签
        label = svg.selectAll('.node-label')
            .data(nodes)
            .enter().append('text')
            .attr('class', 'node-label')
            .text(d => d.label);

        // 更新力导向图
        simulation.nodes(nodes).on('tick', ticked);
        simulation.force('link').links(links);
        simulation.alpha(1).restart();
    }

    function ticked() {
        link
            .attr('x1', d => d.source.x)
            .attr('y1', d => d.source.y)
            .attr('x2', d => d.target.x)
            .attr('y2', d => d.target.y);

        node
            .attr('cx', d => d.x)
            .attr('cy', d => d.y);

        label
            .attr('x', d => d.x)
            .attr('y', d => d.y + 30);
    }

    function dragstarted(event, d) {
        if (!event.active) simulation.alphaTarget(0.3).restart();
        d.fx = d.x;
        d.fy = d.y;
    }

    function dragged(event, d) {
        d.fx = event.x;
        d.fy = event.y;
    }

    function dragended(event, d) {
        if (!event.active) simulation.alphaTarget(0);
        d.fx = null;
        d.fy = null;
    }

    function addLegend() {
        const legend = d3.select('#graphContainer')
            .append('div')
            .attr('class', 'graph-legend');

        const items = [
            { type: 'character', label: '角色', color: '#FFB800' },
            { type: 'location', label: '地点', color: '#5FB878' },
            { type: 'item', label: '物品', color: '#1E9FFF' },
            { type: 'event', label: '事件', color: '#FF5722' }
        ];

        items.forEach(item => {
            legend.append('div')
                .attr('class', 'legend-item')
                .html(`
                    <div class="legend-color" style="background: ${item.color};"></div>
                    <span>${item.label}</span>
                `);
        });
    }

    function getQueryParam(name) {
        const urlParams = new URLSearchParams(window.location.search);
        return urlParams.get(name);
    }
});
