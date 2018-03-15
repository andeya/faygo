<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>template database access demo</title>
</head>
<body>
   {% for o in SimpleData("biz/demo","tpldemo","{}") %}
       <h3>{{o.id}},{{o.nick}},员工名称:{{ o.cnname }}</h3>
   {% endfor %}
   <hr>
    {% for o in SimpleData2("biz/demo","tpldemo","{}") %}
       <h3>{{o.id}},{{o.nick}},员工名称:{{ o.cnname }}</h3>
   {% endfor %}
</body>
</html>
