package cli

import "fmt"

func PrintUsage(choice MenuChoice) {
	switch choice {
	case MenuMain:
		fmt.Println("GitHub Issues Manager CLI")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  issues <action> [arguments] [flags]")
		fmt.Println()
		fmt.Println("Actions:")
		fmt.Println("  create   Create a new issue")
		fmt.Println("  read     Read a single issue by number")
		fmt.Println("  update   Update an existing issue")
		fmt.Println("  delete   Delete/close an issue")
		fmt.Println("  search   Search issues by criteria/query terms")
		fmt.Println("  exit     Exit the program")
	case MenuCreate:
		fmt.Println("Create usage:")
		fmt.Println("  issues create [--title <title>] [--body <body>] [--assignee <login>] [--assignees <a,b>] [--labels <l1,l2>] [--milestone <number|name>] [--type <type>] [title]")
		fmt.Println("  Uses GITHUB_OWNER and GITHUB_REPO from environment")
	case MenuRead:
		fmt.Println("Read usage:")
		fmt.Println("  issues read <issue_number>")
		fmt.Println("  Uses GITHUB_OWNER and GITHUB_REPO from environment")
	case MenuUpdate:
		fmt.Println("Update usage:")
		fmt.Println("  issues update <issue_number> [--title <title>] [--body <body>] [--assignee <login>] [--assignees <a,b>] [--labels <l1,l2>] [--milestone <number|name>] [--type <type>]")
		fmt.Println("  Uses GITHUB_OWNER and GITHUB_REPO from environment")
	case MenuDelete:
		fmt.Println("Delete usage:")
		fmt.Println("  issues delete <issue_number>")
		fmt.Println("  Uses GITHUB_OWNER and GITHUB_REPO from environment")
	case MenuSearch:
		fmt.Println("Search usage:")
		fmt.Println("  issues search <terms...>")
		fmt.Println("  issues search repo:owner/repo is:open label:bug")
	default:
		PrintUsage(MenuMain)
	}
}
