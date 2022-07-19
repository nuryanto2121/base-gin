package templateemail

const (
	VerifyCode = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<title>Title of the document</title>
	</head>
	<body>

	<h4>Hai {Name}</h4>
	
	<h1>{GenerateCode}</h1>


	</body>
	</html>

	`
)
