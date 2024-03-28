package main

var HTML_TEMPLATE = `<!DOCTYPE html>
<html>
<head>
	<title>Welcome to Billowdev</title>
</head>
<body>
	<table width="100%" cellpadding="0" cellspacing="0" style="font-family: Arial, sans-serif;">
		<tr>
			<td align="center" bgcolor="#f4f4f4">
				<table width="600" cellpadding="10" cellspacing="0" style="background-color: #ffffff;">
					<!-- Header -->
					<tr>
						<td bgcolor="#6093E2" align="center" style="color: #ffffff; padding: 10px;">
							<h1>Your Test Email</h1>
						</td>
					</tr>

				<div style="width: max-content; margin: 0 auto; text-align: center">
					<img
					  style="width: 25%; margin-bottom: 5px"
					  src="https://yt3.googleusercontent.com/PFmOY2nzI91G1-J6fSj_uVO0KmspprKjAbvMhA3vznSKMy-n3yG9eFDTSMLSYT0T1ELWJkPTLA=s176-c-k-c0x00ffffff-no-rj"
					/>
				  </div>

					<!-- Content -->
					<tr>
						<td style="padding: 20px;">
							<p>Dear {{.CustomerName}},</p>
							
							<p>Your test email</p>

							<p>Best regards,</p>
							<p>Billowdev Team</p>
						</td>
					</tr>
					<!-- Footer -->
					<tr>
						<td bgcolor="#6093E2" align="center" style="color: #ffffff; padding: 10px;">
							<p>&copy; 2024 Billowdev. All rights reserved.</p>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
</body>
</html>`
