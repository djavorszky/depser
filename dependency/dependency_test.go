package dependency

import (
	"testing"
)

func TestDependency_doAddDependency(t *testing.T) {
	type args struct {
		depender  string
		dependent string
	}
	tests := []struct {
		name string
		args []args
	}{
		{"First", []args{args{"First", "depender"}}},
		{"Two runs", []args{args{"Two runs", "depender1"}, args{"Two runs", "depender2"}}},
		{"Three runs", []args{args{"Three runs", "depender1"}, args{"Three runs", "depender2"}, args{"Three runs", "depender3"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			for _, depTest := range tt.args {
				d.mustAddDependency(depTest.depender, depTest.dependent)

				res, ok := d.deps[depTest.depender]
				if !ok {
					t.Errorf("failed to add dependency: %s -> %s", depTest.depender, depTest.dependent)
				}

				found := false
				for _, d := range res {
					if d == depTest.dependent {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("didn't find added dependency: %s -> %s", depTest.depender, depTest.dependent)
				}
			}

			actualLen := len(d.deps[tt.name])
			if actualLen != len(tt.args) {
				t.Errorf("length of resulting dependency list is not as expected. Expected: %d, actual: %d", len(tt.args), actualLen)
			}
		})
	}
}

func TestDependency_Add(t *testing.T) {
	type args struct {
		depender  string
		dependent string
		wantErr   bool
	}
	tests := []struct {
		name string
		args []args
	}{
		{"First", []args{args{"First", "dependee", false}}},
		{"Two runs", []args{args{"Two runs", "dependee1", false}, args{"Two runs", "dependee2", false}}},
		{"Three runs", []args{args{"Three runs", "dependee1", false}, args{"Three runs", "dependee2", false}, args{"Three runs", "dependee3", false}}},

		{"Missing dependent", []args{args{"", "dependee", true}}},
		{"Missing dependee", []args{args{"First", "", true}}},
		{"Dependency cycle", []args{args{"a", "b", false}, args{"b", "a", true}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			for _, depTest := range tt.args {
				if err := d.Add(depTest.depender, depTest.dependent); (err != nil) != depTest.wantErr {
					t.Errorf("Dependency.Add() error = %v, wantErr %v", err, depTest.wantErr)
					return
				}

				if depTest.wantErr {
					return
				}
			}

			// one depender, many dependents
			depLen := len(d.deps[tt.name])
			if len(tt.args) != depLen {
				t.Errorf("length of resulting dependency list is not as expected. Expected: %d, actual: %d", len(tt.args), depLen)
			}

			// many stalked, only 1 stalker each
			visLen := len(d.visibilities)
			if len(tt.args) != visLen {
				t.Errorf("length of resulting visibility list is not as expected. Expected: %d, actual: %d", len(tt.args), visLen)
			}
		})
	}
}

func TestDependency_mustAddVisibility(t *testing.T) {
	type args struct {
		stalked string
		stalker string
	}
	tests := []struct {
		name string
		args []args
	}{
		{"First", []args{args{"First", "stalker"}}},
		{"Two runs", []args{args{"Two runs", "stalker1"}, args{"Two runs", "stalker2"}}},
		{"Three runs", []args{args{"Three runs", "stalker1"}, args{"Three runs", "stalker2"}, args{"Three runs", "stalker3"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			for _, depTest := range tt.args {
				d.mustAddVisibility(depTest.stalked, depTest.stalker)

				res, ok := d.visibilities[depTest.stalked]
				if !ok {
					t.Errorf("failed to add visibility: %s -> %s", depTest.stalked, depTest.stalker)
					continue
				}

				found := false
				for _, d := range res {
					if d == depTest.stalker {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("didn't find added visibility: %s -> %s", depTest.stalked, depTest.stalker)
				}
			}

			actualLen := len(d.visibilities[tt.name])
			if actualLen != len(tt.args) {
				t.Errorf("length of resulting visibility list is not as expected. Expected: %d, actual: %d", len(tt.args), actualLen)
			}
		})
	}
}

func TestDependency_checkCycles(t *testing.T) {
	type args struct {
		depender  string
		dependent string
	}
	tests := []struct {
		name     string
		depender string
		args     []args
		wantErr  bool
	}{
		{"Works", "a", []args{args{"a", "b"}, args{"b", "c"}}, false},

		{"Direct cycle", "a", []args{args{"a", "b"}, args{"b", "a"}}, true},
		{"Indirect cycle", "a", []args{args{"a", "b"}, args{"b", "c"}, args{"c", "a"}}, true},
		{"Long cycle", "a", []args{args{"a", "b"}, args{"b", "c"}, args{"c", "d"}, args{"d", "e"}, args{"e", "f"}, args{"f", "a"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			for _, arg := range tt.args {
				d.mustAddDependency(arg.depender, arg.dependent)
			}
			if err := d.checkCycles(tt.depender); (err != nil) != tt.wantErr {
				t.Errorf("Dependency.checkCycles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_trimToCycle(t *testing.T) {
	type args struct {
		cycle    string
		offender string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"double",
			args{
				"com.liferay.portal.tools.theme.builder.ThemeBuilder -> com.liferay.portal.tools.theme.builder.internal.util.FileUtil -> com.liferay.portal.tools.theme.builder.ThemeBuilder -> com.liferay.portal.tools.theme.builder.internal.util.FileUtil",
				"com.liferay.portal.tools.theme.builder.internal.util.FileUtil",
			},
			"com.liferay.portal.tools.theme.builder.internal.util.FileUtil -> com.liferay.portal.tools.theme.builder.ThemeBuilder -> com.liferay.portal.tools.theme.builder.internal.util.FileUtil",
			false,
		},
		{
			"double_2",
			args{
				"com.liferay.portal.tools.theme.builder.ThemeBuilder -> com.liferay.portal.tools.theme.builder.internal.util.FileUtil -> com.liferay.portal.tools.theme.builder.ThemeBuilder -> com.liferay.portal.tools.theme.builder.internal.util.FileUtil",
				"com.liferay.portal.tools.theme.builder.ThemeBuilder",
			},
			"com.liferay.portal.tools.theme.builder.ThemeBuilder -> com.liferay.portal.tools.theme.builder.internal.util.FileUtil -> com.liferay.portal.tools.theme.builder.ThemeBuilder",
			false,
		},
		{
			"long",
			args{
				"com.liferay.adaptive.media.document.library.web.internal.optimizer.test.DLAMImageOptimizerTest -> com.liferay.portal.test.rule.LiferayIntegrationTestRule -> com.liferay.portal.test.rule.callback.MainServletTestCallback -> com.liferay.portal.servlet.MainServlet -> com.liferay.portal.events.StartupAction -> com.liferay.portal.kernel.resiliency.mpi.MPIHelperUtil -> com.liferay.portal.kernel.resiliency.spi.SPI -> com.liferay.portal.kernel.resiliency.spi.agent.SPIAgent -> com.liferay.portal.kernel.resiliency.spi.SPI",
				"com.liferay.portal.kernel.resiliency.spi.SPI",
			},
			"com.liferay.portal.kernel.resiliency.spi.SPI -> com.liferay.portal.kernel.resiliency.spi.agent.SPIAgent -> com.liferay.portal.kernel.resiliency.spi.SPI",
			false,
		},
		{
			"special",
			args{
				"${package}.content.targeting.report.__className__Report -> com.liferay.portal.kernel.util.ParamUtil -> com.liferay.portal.kernel.service.ServiceContext -> com.liferay.portal.kernel.util.PortalUtil -> com.liferay.expando.kernel.model.ExpandoBridge -> com.liferay.portal.kernel.service.ServiceContext",
				"com.liferay.portal.kernel.service.ServiceContext",
			},
			"com.liferay.portal.kernel.service.ServiceContext -> com.liferay.portal.kernel.util.PortalUtil -> com.liferay.expando.kernel.model.ExpandoBridge -> com.liferay.portal.kernel.service.ServiceContext",
			false,
		},
		{
			"special_long",
			args{
				"${package}.application.list.__className__PanelCategory -> com.liferay.application.list.BasePanelCategory -> com.liferay.application.list.display.context.logic.PanelCategoryHelper -> com.liferay.application.list.PanelAppRegistry -> com.liferay.portal.kernel.portlet.PortletPreferencesFactory -> com.liferay.portal.kernel.theme.ThemeDisplay -> com.liferay.exportimport.kernel.staging.StagingUtil -> com.liferay.exportimport.kernel.lar.PortletDataContext -> com.liferay.portal.kernel.model.ClassedModel -> com.liferay.expando.kernel.model.ExpandoBridge -> com.liferay.portal.kernel.service.ServiceContext -> com.liferay.portal.kernel.portlet.PortletPreferencesFactoryUtil -> com.liferay.portal.kernel.theme.ThemeDisplay",
				"com.liferay.portal.kernel.theme.ThemeDisplay",
			},
			"com.liferay.portal.kernel.theme.ThemeDisplay -> com.liferay.exportimport.kernel.staging.StagingUtil -> com.liferay.exportimport.kernel.lar.PortletDataContext -> com.liferay.portal.kernel.model.ClassedModel -> com.liferay.expando.kernel.model.ExpandoBridge -> com.liferay.portal.kernel.service.ServiceContext -> com.liferay.portal.kernel.portlet.PortletPreferencesFactoryUtil -> com.liferay.portal.kernel.theme.ThemeDisplay",
			false,
		},

		{
			"nothing_to_do",
			args{
				"com.liferay.portal.tools.theme.builder.ThemeBuilder -> com.liferay.portal.tools.theme.builder.internal.util.FileUtil -> com.liferay.portal.tools.theme.builder.ThemeBuilder",
				"com.liferay.portal.tools.theme.builder.ThemeBuilder",
			},
			"com.liferay.portal.tools.theme.builder.ThemeBuilder -> com.liferay.portal.tools.theme.builder.internal.util.FileUtil -> com.liferay.portal.tools.theme.builder.ThemeBuilder",
			false,
		},

		{"missing_cycle", args{"", "com.liferay.portal.kernel.theme.ThemeDisplay"}, "", true},

		{"missing_offender", args{"com.liferay.portal.kernel.service.ServiceContext -> com.liferay.portal.kernel.util.PortalUtil -> com.liferay.expando.kernel.model.ExpandoBridge -> com.liferay.portal.kernel.service.ServiceContext", ""}, "", true},

		{"offender_not_found", args{"com.liferay.portal.tools.theme.builder.ThemeBuilder -> com.liferay.portal.tools.theme.builder.internal.util.FileUtil -> com.liferay.portal.tools.theme.builder.ThemeBuilder", "com.liferay.portal.kernel.theme.ThemeDisplay"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := trimToCycle(tt.args.cycle, tt.args.offender)
			if (err != nil) != tt.wantErr {
				t.Errorf("trimToCycle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("trimToCycle() = %v, want %v", got, tt.want)
			}
		})
	}
}
