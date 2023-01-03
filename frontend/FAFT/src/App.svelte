<script>
  import ShowResult from "./FetchData.svelte";
  import { Router, Route, Link } from "svelte-navigator";
  import ReadMine from "./ReadMine.svelte";
  import Working from "./WorkRelevent.svelte";
  import Stainless from "./Stainless.svelte";
</script>

<head>
  <title>FAFT</title>
</head>
<Router>
  <div class="container">
    <nav>
      <sp>
        <Link to="home">Home</Link>
        <span>&emsp;&emsp;&emsp;&emsp;</span>
        <Link to="stainless">Stainless</Link>
        <span>&emsp;&emsp;&emsp;&emsp;</span>
        <Link to="result">Test Ressult</Link>
        <span>&emsp;&emsp;&emsp;&emsp;</span>
        <Link to="readmineList">RedMine</Link>
        <span>&emsp;&emsp;&emsp;&emsp;</span>
        <Link to="info">Knowledge</Link>
      </sp>
    </nav>
  </div>

  <main>
    <Route path="home">
      <h1 style="text-align: center; font-size: 60px">FAFT</h1>
      <Working />
    </Route>
    <Route path="stainless">
      <Stainless />
    </Route>
    <Route path="result">
      <h1>FAFT</h1>
      <ShowResult />
    </Route>

    <Route path="readmineList">
      <h1>Kevin</h1>
      <ReadMine />
    </Route>
  </main>
</Router>

<style>
  .container {
    margin-top: 5px;
    margin-left: 40%;
    margin-right: auto;
    margin-bottom: 20px;
    border-style: outset;
    border-width: 9px;
    font-style: italic;
    text-align: right;
  }
  sp {
    font-size: 20px;
  }
</style>
