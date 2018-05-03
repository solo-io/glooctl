package route

import (
	"fmt"

	"github.com/solo-io/gloo/pkg/bootstrap/configstorage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	helper "github.com/solo-io/glooctl/internal/test-helper"
)

var _ = Describe("Creating route", func() {
	BeforeEach(helper.SetupStorage)
	AfterEach(helper.TearDownStorage)

	It("should allow creating without any virtual service in system", func() {
		helper.RunWithArgs("route", "create", "--path-prefix", "/foo", "--upstream", "test-upstream").
			ExpectExitCodeAndOutput(0, "Did not find a default virtual service. Creating",
				"/foo", "test-upstream")
	})

	It("should allow creating with a default virtual service not called 'default'", func() {
		// create a default virtual service of different name - mydefault
		fmt.Fprintln(GinkgoWriter, "Creating a virtual host with different name")
		helper.RunWithArgs("virtualservice", "create", "-f", "testdata/vs-mydefault.yaml").
			ExpectExitCode(0)

		fmt.Fprintln(GinkgoWriter, "Creating a route")
		helper.RunWithArgs("route", "create", "--path-prefix", "/foo", "--upstream", "test-upstream").
			ExpectExitCodeAndOutput(0, "Using virtual service: mydefault", "/foo", "test-upstream")
	})

	It("should allow selecting the virtual service using domain", func() {
		// create a default virtual service of different name - mydefault
		helper.RunWithArgs("virtualservice", "create", "-f", "testdata/vs-mydefault.yaml").
			ExpectExitCode(0)
		// create a virtual service with domain
		helper.RunWithArgs("virtualservice", "create", "-f", "testdata/vs-with-domain.yaml").
			ExpectExitCode(0)

		fmt.Fprintln(GinkgoWriter, "Creating a route with domain")
		helper.RunWithArgs("route", "create", "--domain", "axhixh.com", "--path-prefix", "/foo",
			"--upstream", "test-upstream").
			ExpectExitCodeAndOutput(0, "Using virtual service: with-domain", "/foo", "test-upstream")
	})

	It("should fail when giving invalid domain", func() {
		helper.RunWithArgs("route", "create", "--domain", "nowhere.com", "--path-prefix", "/foo",
			"--upstream", "test-upstream").
			ExpectExitCodeAndOutput(1, "didn't find any virtual service for the domain nowhere.com")
	})

	It("should reorder the routes if asked to sort", func() {
		helper.RunWithArgs("route", "create", "--path-prefix", "/foo", "--upstream", "test-upstream").
			ExpectExitCode(0)

		// create second with sort
		helper.RunWithArgs("route", "create", "--path-exact", "/a", "--upstream", "test-upstream2",
			"--sort").ExpectExitCode(0)

		sc, err := configstorage.Bootstrap(*helper.BootstrapOpts())
		Expect(err).NotTo(HaveOccurred())
		vs, err := sc.V1().VirtualServices().Get("default")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(vs.Routes)).To(Equal(2))
		matcher, ok := vs.Routes[0].Matcher.(*v1.Route_RequestMatcher)
		Expect(ok).To(Equal(true))
		path, ok := matcher.RequestMatcher.Path.(*v1.RequestMatcher_PathExact)
		Expect(ok).To(Equal(true))
		Expect(path.PathExact).To(Equal("/a"))
	})
})
