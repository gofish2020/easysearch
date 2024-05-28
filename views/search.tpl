<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{ .keyword }} - {{ .title }}</title>
  <link rel="stylesheet" href="https://cdn.bootcdn.net/ajax/libs/twitter-bootstrap/5.2.1/css/bootstrap.min.css">
  <script src="https://cdn.bootcdn.net/ajax/libs/vue/2.6.1/vue.js"></script>
  <script src="https://cdn.bootcdn.net/ajax/libs/jquery/3.6.1/jquery.min.js"></script>
  <script src="https://cdn.bootcdn.net/ajax/libs/twitter-bootstrap/5.2.1/js/bootstrap.min.js"></script>
  <script src="https://cdn.bootcdn.net/ajax/libs/echarts/5.5.0/echarts.min.js"></script>
</head>
<body>
  <div id="title">{{ .title }}</div>
  <div id="time">
    查询耗时：{{ .latency }}
    <p>返回时间:{{ .time }}</p>
  </div>
  <div id="content">
    <form action="/" method="get">
      <div class="row">
        <div class="col-md-8">
          <input type="text" name="keyword" class="form-control" id="keyword" value="{{ .keyword }}">
        </div>
        <div class="col-md-2 col-md-offset-1">
          <button class="btn btn-success form-control" type="submit">Go搜搜Go</button>
        </div>
      </div>
    </form>
    <div id="result" class="row">
      <ul>
        {{range .values}}
          <li>
            <h4><a target="_blank" href="{{ .Url }}">{{ .Title }}</a></h4>
            <p>{{ .Brief }}</p>
          </li>
        {{end}}
      </ul>
    </div>
  </div>
<style>
* {
  margin: 0;
  padding: 0;
}
#title {
  font-size: 32px;
  text-align: center;
  width: 100%;
  margin-top: 20px;
  color: green;
}
#time {
  font-size: 16px;
  margin: 10px;
  text-align: center;
  line-height: 2em;
}
#content {
  width: 600px;
  margin: 20px auto;
}
#result {
  margin-top: 50px;
}
ul, li {
  list-style: none;
}
</style>
</body>
</html>