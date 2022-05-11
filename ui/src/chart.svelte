<script>
    import Chart from 'chart.js/auto';
	import { onMount } from 'svelte';

	export let jsonData = {}
	let chartCanvas
	let labels = []
	let like_data = []
	let retweet_data = []
	let reply_data = []
	let quote_data = []
	let chart

	function resetData () {
		chart.reset()
		labels = []
		like_data = []
		retweet_data = []
		reply_data = []
		quote_data = []
	}

	const data = {
		labels: labels,
		datasets: [{
			label: 'Likes',
			data: like_data,
			fill: false,
			borderColor: 'rgb(75, 192, 192)',
			tension: 0.1
		},
		{
			label: 'Retweets',
			data: retweet_data,
			fill: false,
			borderColor: 'rgb(255, 99, 132)',
			tension: 0.1
		},
		{
			label: 'Quotes retweets',
			data: quote_data,
			fill: false,
			borderColor: 'rgb(255, 159, 64)',
			tension: 0.1
		},
		{
			label: 'Replies',
			data: reply_data,
			fill: false,
			borderColor: 'rgb(255, 205, 86)',
			tension: 0.1
		}
	]
	};

	onMount(() => {
		const ctx = chartCanvas.getContext("2d");
		chart = new Chart(ctx, {
			type: 'line',
			data: data,
			options: {
				maintainAspectRatio: false,
				responsive: true,
			}
		});
        console.log(jsonData)
        console.log(labels)
        console.log("Log me!")
        for (const key of Object.keys(jsonData)) {
            labels = labels.concat(key)
            like_data = like_data.concat(jsonData[key]["like_count"])
            retweet_data = retweet_data.concat(jsonData[key]["retweet_count"])
            reply_data = reply_data.concat(jsonData[key]["reply_count"])
            quote_data = quote_data.concat(jsonData[key]["quote_count"])
        }
        console.log(labels)
        chart.data.labels = labels
        chart.data.datasets[0].data = like_data
        chart.data.datasets[1].data = retweet_data
        chart.data.datasets[2].data = quote_data
        chart.data.datasets[3].data = reply_data
        chart.update()
	});
</script>

<canvas bind:this={chartCanvas}></canvas>

<style>
    canvas {
        margin-bottom: 40px;
    }
</style>
