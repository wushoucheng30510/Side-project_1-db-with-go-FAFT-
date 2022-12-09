<script>
  function fetchData() {
    return fetch("http://10.240.102.12:8082/test").then((res) => res.json());
  }
</script>

<table style="table-layout:fixed;">
  <thead>
    <tr>
      <th colspan="10" style="background-color:blue;">Result</th>
    </tr>
    <tr>
      <td style="text-align:center">ID</td>
      <td style="text-align:center">Time</td>
      <td style="text-align:center">Tester</td>
      <td style="text-align:center">Name</td>
      <td style="text-align:center">Board</td>
      <td style="text-align:center">Model</td>
      <td style="text-align:center">Version</td>
      <td style="text-align:center">Log_Path</td>
      <td style="text-align:center">Result</td>
      <td style="text-align:center">Reason</td>
    </tr>
  </thead>
  <tbody>
    {#await fetchData()}
      <span>Loading...</span>
    {:then postTest}
      {#each postTest as post (post.id)}
        <tr>
          <td class="blue">{post.id}</td>
          <td class="time">{post.time}</td>
          <td style="text-align:center">{post.tester}</td>
          <td style="text-align:center">{post.name}</td>
          <td style="text-align:center">{post.board}</td>
          <td style="text-align:center">{post.model}</td>
          <td>{post.version}</td>
          <td>
            <a href="http://10.240.102.12:8082/log/{post.logPath}"
              >{post.logPath}</a
            >
          </td>
          {#if post.passOrFail == "Pass"}
            <td class="green" style="text-align:center">{post.passOrFail}</td>
          {:else}
            <td class="red" style="text-align:center">{post.passOrFail}</td>
          {/if}
          <td class="red" style="table-layout:fixed;">{post.reason}</td>
        </tr>
      {/each}
    {:catch error}
      <span>{error}</span>
    {/await}
  </tbody>
</table>

<style>
  .red {
    color: red;
    font-size: 16px;
  }

  table,
  td {
    border: 1px solid #333;
    font-size: 16px;
  }

  thead {
    background-color: #333;
    color: #fff;
  }

  .blue {
    color: blue;
  }

  .green {
    color: rgb(28, 118, 1);
  }
  .time {
    font-size: 10px;
  }
</style>
