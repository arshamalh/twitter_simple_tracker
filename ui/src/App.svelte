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
				resolve(JSON.parse(e.target.result))
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
				<p>
					<a href="https://twitter.com/i/user/{jsonData['author_id']}">Author</a> - <a href="https://twitter.com/i/web/status/{jsonData['id']}">Tweet</a> - {jsonData["created_at"]}
				</p>
				<Chart jsonData={jsonData.public_metrics}/>
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

	p {
		text-align: center;
	}

	input {
		margin-bottom: 10px;
	}
</style>