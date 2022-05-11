<script>
	import Chart from './chart.svelte';

	let files = []

	// function setData(event) {
	// 	const fileReader = new FileReader();
	// 	for (const file of event.target?.files) {
	// 		fileReader.readAsText(file)
	// 		fileReader.onload = (e) => {
	// 			jsonData = JSON.parse(e.target.result).public_metrics
	// 			chart.data = jsonData
	// 		}
	// 	}
	// }

	function extractJSONDataFromFile(file) {
		return new Promise((resolve, reject) => {
			const fileReader = new FileReader();
			fileReader.readAsText(file)
			fileReader.onload = (e) => {
				resolve(JSON.parse(e.target.result).public_metrics)
			}
		})
	}
</script>

<main>
	<input type="file" multiple bind:files>
	{#each files as file}
		<div>
			{#await extractJSONDataFromFile(file)}
				<p>waiting...</p>
			{:then jsonData}
				<h3>{file.name}</h3>
				<Chart jsonData={jsonData}/>
			{/await}
		</div>
	{/each}
</main>

<style>
	main {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}

	div {
		width: 100%;
		height: 100vh;
		margin-top: 20px;
	}

	h3, p {
		text-align: center;
	}

	input {
		margin-bottom: 10px;
	}
</style>