//go:build windows

package dhcp

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus-community/windows_exporter/internal/mi"
	"github.com/prometheus-community/windows_exporter/internal/perfdata"
	"github.com/prometheus-community/windows_exporter/internal/perfdata/perftypes"
	v1 "github.com/prometheus-community/windows_exporter/internal/perfdata/v1"
	"github.com/prometheus-community/windows_exporter/internal/types"
	"github.com/prometheus-community/windows_exporter/internal/utils"
	"github.com/prometheus/client_golang/prometheus"
)

const Name = "dhcp"

type Config struct{}

var ConfigDefaults = Config{}

// A Collector is a Prometheus Collector perflib DHCP metrics.
type Collector struct {
	config Config

	perfDataCollector perfdata.Collector

	acksTotal                                        *prometheus.Desc
	activeQueueLength                                *prometheus.Desc
	conflictCheckQueueLength                         *prometheus.Desc
	declinesTotal                                    *prometheus.Desc
	deniedDueToMatch                                 *prometheus.Desc
	deniedDueToNonMatch                              *prometheus.Desc
	discoversTotal                                   *prometheus.Desc
	duplicatesDroppedTotal                           *prometheus.Desc
	failoverBndAckReceivedTotal                      *prometheus.Desc
	failoverBndAckSentTotal                          *prometheus.Desc
	failoverBndUpdDropped                            *prometheus.Desc
	failoverBndUpdPendingOutboundQueue               *prometheus.Desc
	failoverBndUpdReceivedTotal                      *prometheus.Desc
	failoverBndUpdSentTotal                          *prometheus.Desc
	failoverTransitionsCommunicationInterruptedState *prometheus.Desc
	failoverTransitionsPartnerDownState              *prometheus.Desc
	failoverTransitionsRecoverState                  *prometheus.Desc
	informsTotal                                     *prometheus.Desc
	nACKsTotal                                       *prometheus.Desc
	offerQueueLength                                 *prometheus.Desc
	offersTotal                                      *prometheus.Desc
	packetsExpiredTotal                              *prometheus.Desc
	packetsReceivedTotal                             *prometheus.Desc
	releasesTotal                                    *prometheus.Desc
	requestsTotal                                    *prometheus.Desc
}

func New(config *Config) *Collector {
	if config == nil {
		config = &ConfigDefaults
	}

	c := &Collector{
		config: *config,
	}

	return c
}

func NewWithFlags(_ *kingpin.Application) *Collector {
	return &Collector{}
}

func (c *Collector) GetName() string {
	return Name
}

func (c *Collector) GetPerfCounter(_ *slog.Logger) ([]string, error) {
	if utils.PDHEnabled() {
		return []string{}, nil
	}

	return []string{"DHCP Server"}, nil
}

func (c *Collector) Close(_ *slog.Logger) error {
	return nil
}

func (c *Collector) Build(_ *slog.Logger, _ *mi.Session) error {
	if utils.PDHEnabled() {
		counters := []string{
			acksTotal,
			activeQueueLength,
			conflictCheckQueueLength,
			declinesTotal,
			deniedDueToMatch,
			deniedDueToNonMatch,
			discoversTotal,
			duplicatesDroppedTotal,
			failoverBndAckReceivedTotal,
			failoverBndAckSentTotal,
			failoverBndUpdDropped,
			failoverBndUpdPendingOutboundQueue,
			failoverBndUpdReceivedTotal,
			failoverBndUpdSentTotal,
			failoverTransitionsCommunicationInterruptedState,
			failoverTransitionsPartnerDownState,
			failoverTransitionsRecoverState,
			informsTotal,
			nacksTotal,
			offerQueueLength,
			offersTotal,
			packetsExpiredTotal,
			packetsReceivedTotal,
			releasesTotal,
			requestsTotal,
		}

		var err error

		c.perfDataCollector, err = perfdata.NewCollector(perfdata.V1, "DHCP Server", perfdata.AllInstances, counters)
		if err != nil {
			return fmt.Errorf("failed to create DHCP Server collector: %w", err)
		}
	}

	c.packetsReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "packets_received_total"),
		"Total number of packets received by the DHCP server (PacketsReceivedTotal)",
		nil,
		nil,
	)
	c.duplicatesDroppedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "duplicates_dropped_total"),
		"Total number of duplicate packets received by the DHCP server (DuplicatesDroppedTotal)",
		nil,
		nil,
	)
	c.packetsExpiredTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "packets_expired_total"),
		"Total number of packets expired in the DHCP server message queue (PacketsExpiredTotal)",
		nil,
		nil,
	)
	c.activeQueueLength = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "active_queue_length"),
		"Number of packets in the processing queue of the DHCP server (ActiveQueueLength)",
		nil,
		nil,
	)
	c.conflictCheckQueueLength = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "conflict_check_queue_length"),
		"Number of packets in the DHCP server queue waiting on conflict detection (ping). (ConflictCheckQueueLength)",
		nil,
		nil,
	)
	c.discoversTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "discovers_total"),
		"Total DHCP Discovers received by the DHCP server (DiscoversTotal)",
		nil,
		nil,
	)
	c.offersTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "offers_total"),
		"Total DHCP Offers sent by the DHCP server (OffersTotal)",
		nil,
		nil,
	)
	c.requestsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "requests_total"),
		"Total DHCP Requests received by the DHCP server (RequestsTotal)",
		nil,
		nil,
	)
	c.informsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "informs_total"),
		"Total DHCP Informs received by the DHCP server (InformsTotal)",
		nil,
		nil,
	)
	c.acksTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "acks_total"),
		"Total DHCP Acks sent by the DHCP server (AcksTotal)",
		nil,
		nil,
	)
	c.nACKsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "nacks_total"),
		"Total DHCP Nacks sent by the DHCP server (NacksTotal)",
		nil,
		nil,
	)
	c.declinesTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "declines_total"),
		"Total DHCP Declines received by the DHCP server (DeclinesTotal)",
		nil,
		nil,
	)
	c.releasesTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "releases_total"),
		"Total DHCP Releases received by the DHCP server (ReleasesTotal)",
		nil,
		nil,
	)
	c.offerQueueLength = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "offer_queue_length"),
		"Number of packets in the offer queue of the DHCP server (OfferQueueLength)",
		nil,
		nil,
	)
	c.deniedDueToMatch = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "denied_due_to_match_total"),
		"Total number of DHCP requests denied, based on matches from the Deny list (DeniedDueToMatch)",
		nil,
		nil,
	)
	c.deniedDueToNonMatch = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "denied_due_to_nonmatch_total"),
		"Total number of DHCP requests denied, based on non-matches from the Allow list (DeniedDueToNonMatch)",
		nil,
		nil,
	)
	c.failoverBndUpdSentTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "failover_bndupd_sent_total"),
		"Number of DHCP fail over Binding Update messages sent (FailoverBndupdSentTotal)",
		nil,
		nil,
	)
	c.failoverBndUpdReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "failover_bndupd_received_total"),
		"Number of DHCP fail over Binding Update messages received (FailoverBndupdReceivedTotal)",
		nil,
		nil,
	)
	c.failoverBndAckSentTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "failover_bndack_sent_total"),
		"Number of DHCP fail over Binding Ack messages sent (FailoverBndackSentTotal)",
		nil,
		nil,
	)
	c.failoverBndAckReceivedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "failover_bndack_received_total"),
		"Number of DHCP fail over Binding Ack messages received (FailoverBndackReceivedTotal)",
		nil,
		nil,
	)
	c.failoverBndUpdPendingOutboundQueue = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "failover_bndupd_pending_in_outbound_queue"),
		"Number of pending outbound DHCP fail over Binding Update messages (FailoverBndupdPendingOutboundQueue)",
		nil,
		nil,
	)
	c.failoverTransitionsCommunicationInterruptedState = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "failover_transitions_communicationinterrupted_state_total"),
		"Total number of transitions into COMMUNICATION INTERRUPTED state (FailoverTransitionsCommunicationinterruptedState)",
		nil,
		nil,
	)
	c.failoverTransitionsPartnerDownState = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "failover_transitions_partnerdown_state_total"),
		"Total number of transitions into PARTNER DOWN state (FailoverTransitionsPartnerdownState)",
		nil,
		nil,
	)
	c.failoverTransitionsRecoverState = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "failover_transitions_recover_total"),
		"Total number of transitions into RECOVER state (FailoverTransitionsRecoverState)",
		nil,
		nil,
	)
	c.failoverBndUpdDropped = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "failover_bndupd_dropped_total"),
		"Total number of DHCP fail over Binding Updates dropped (FailoverBndupdDropped)",
		nil,
		nil,
	)

	return nil
}

func (c *Collector) Collect(ctx *types.ScrapeContext, logger *slog.Logger, ch chan<- prometheus.Metric) error {
	if utils.PDHEnabled() {
		return c.collectPDH(ch)
	}

	logger = logger.With(slog.String("collector", Name))

	return c.collect(ctx, logger, ch)
}

func (c *Collector) collect(ctx *types.ScrapeContext, logger *slog.Logger, ch chan<- prometheus.Metric) error {
	var dhcpPerfs []dhcpPerf

	if err := v1.UnmarshalObject(ctx.PerfObjects["DHCP Server"], &dhcpPerfs, logger); err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.packetsReceivedTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].PacketsReceivedTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.duplicatesDroppedTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].DuplicatesDroppedTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.packetsExpiredTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].PacketsExpiredTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.activeQueueLength,
		prometheus.GaugeValue,
		dhcpPerfs[0].ActiveQueueLength,
	)

	ch <- prometheus.MustNewConstMetric(
		c.conflictCheckQueueLength,
		prometheus.GaugeValue,
		dhcpPerfs[0].ConflictCheckQueueLength,
	)

	ch <- prometheus.MustNewConstMetric(
		c.discoversTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].DiscoversTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.offersTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].OffersTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.requestsTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].RequestsTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.informsTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].InformsTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.acksTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].AcksTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.nACKsTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].NacksTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.declinesTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].DeclinesTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.releasesTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].ReleasesTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.offerQueueLength,
		prometheus.GaugeValue,
		dhcpPerfs[0].OfferQueueLength,
	)

	ch <- prometheus.MustNewConstMetric(
		c.deniedDueToMatch,
		prometheus.CounterValue,
		dhcpPerfs[0].DeniedDueToMatch,
	)

	ch <- prometheus.MustNewConstMetric(
		c.deniedDueToNonMatch,
		prometheus.CounterValue,
		dhcpPerfs[0].DeniedDueToNonMatch,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverBndUpdSentTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].FailoverBndUpdSentTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverBndUpdReceivedTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].FailoverBndUpdReceivedTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverBndAckSentTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].FailoverBndAckSentTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverBndAckReceivedTotal,
		prometheus.CounterValue,
		dhcpPerfs[0].FailoverBndAckReceivedTotal,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverBndUpdPendingOutboundQueue,
		prometheus.GaugeValue,
		dhcpPerfs[0].FailoverBndUpdPendingOutboundQueue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverTransitionsCommunicationInterruptedState,
		prometheus.CounterValue,
		dhcpPerfs[0].FailoverTransitionsCommunicationInterruptedState,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverTransitionsPartnerDownState,
		prometheus.CounterValue,
		dhcpPerfs[0].FailoverTransitionsPartnerDownState,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverTransitionsRecoverState,
		prometheus.CounterValue,
		dhcpPerfs[0].FailoverTransitionsRecoverState,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverBndUpdDropped,
		prometheus.CounterValue,
		dhcpPerfs[0].FailoverBndUpdDropped,
	)

	return nil
}

func (c *Collector) collectPDH(ch chan<- prometheus.Metric) error {
	perfData, err := c.perfDataCollector.Collect()
	if err != nil {
		return fmt.Errorf("failed to collect DHCP Server metrics: %w", err)
	}

	data, ok := perfData[perftypes.EmptyInstance]
	if !ok {
		return errors.New("perflib query for DHCP Server returned empty result set")
	}

	ch <- prometheus.MustNewConstMetric(
		c.packetsReceivedTotal,
		prometheus.CounterValue,
		data[packetsReceivedTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.duplicatesDroppedTotal,
		prometheus.CounterValue,
		data[duplicatesDroppedTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.packetsExpiredTotal,
		prometheus.CounterValue,
		data[packetsExpiredTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.activeQueueLength,
		prometheus.GaugeValue,
		data[activeQueueLength].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.conflictCheckQueueLength,
		prometheus.GaugeValue,
		data[conflictCheckQueueLength].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.discoversTotal,
		prometheus.CounterValue,
		data[discoversTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.offersTotal,
		prometheus.CounterValue,
		data[offersTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.requestsTotal,
		prometheus.CounterValue,
		data[requestsTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.informsTotal,
		prometheus.CounterValue,
		data[informsTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.acksTotal,
		prometheus.CounterValue,
		data[acksTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.nACKsTotal,
		prometheus.CounterValue,
		data[nacksTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.declinesTotal,
		prometheus.CounterValue,
		data[declinesTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.releasesTotal,
		prometheus.CounterValue,
		data[releasesTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.offerQueueLength,
		prometheus.GaugeValue,
		data[offerQueueLength].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.deniedDueToMatch,
		prometheus.CounterValue,
		data[deniedDueToMatch].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.deniedDueToNonMatch,
		prometheus.CounterValue,
		data[deniedDueToNonMatch].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverBndUpdSentTotal,
		prometheus.CounterValue,
		data[failoverBndUpdSentTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverBndUpdReceivedTotal,
		prometheus.CounterValue,
		data[failoverBndUpdReceivedTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverBndAckSentTotal,
		prometheus.CounterValue,
		data[failoverBndAckSentTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverBndAckReceivedTotal,
		prometheus.CounterValue,
		data[failoverBndAckReceivedTotal].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverBndUpdPendingOutboundQueue,
		prometheus.GaugeValue,
		data[failoverBndUpdPendingOutboundQueue].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverTransitionsCommunicationInterruptedState,
		prometheus.CounterValue,
		data[failoverTransitionsCommunicationInterruptedState].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverTransitionsPartnerDownState,
		prometheus.CounterValue,
		data[failoverTransitionsPartnerDownState].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverTransitionsRecoverState,
		prometheus.CounterValue,
		data[failoverTransitionsRecoverState].FirstValue,
	)

	ch <- prometheus.MustNewConstMetric(
		c.failoverBndUpdDropped,
		prometheus.CounterValue,
		data[failoverBndUpdDropped].FirstValue,
	)

	return nil
}
